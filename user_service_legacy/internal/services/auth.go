package services

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user_service/config"
	"user_service/internal/models"
	"user_service/internal/repositories"
)

type Auth interface {
	SigIn(user *models.User) (error, *models.TokenDetails)
	SignUp(user *models.User, roleId uint64) (error, *models.TokenDetails)
	Logout(r *http.Request) error
	RefreshToken(refreshToken string) (error, *models.TokenDetails)
	CheckAuthorization(r *http.Request) (error, uint64)
}

type AuthService struct {
	repo   repositories.Auth
	config *config.Config
}

func NewAuthService(db repositories.Auth, cfg *config.Config) *AuthService {
	return &AuthService{
		repo:   db,
		config: cfg,
	}
}

func (as *AuthService) SigIn(user *models.User) (error, *models.TokenDetails) {
	err, _ := as.repo.IsUserExist(as.hashPassword(user.PasswordHash), user.Email)
	if err != nil {
		return err, nil
	}
	return as.createToken(user.ID)
}

func (as *AuthService) SignUp(user *models.User, roleId uint64) (error, *models.TokenDetails) {
	user.PasswordHash = as.hashPassword(user.PasswordHash)
	err, id := as.repo.CreateUser(user, roleId)
	if err != nil {
		return err, nil
	}
	return as.createToken(id)
}

func (as *AuthService) hashPassword(password string) string {
	sha := sha1.New()
	sha.Write([]byte(password))
	sha.Write([]byte(as.config.UserPasswordHashSalt))

	return fmt.Sprintf("%x", sha.Sum(nil))
}

func (as *AuthService) Logout(r *http.Request) error {
	tokenAuth, err := as.extractTokenMetadata(r)
	if err != nil {
		return err
	}
	err = as.repo.DeleteToken(tokenAuth.AccessUuid)
	if err != nil {
		return err
	}
	return nil
}

func (as *AuthService) RefreshToken(refreshToken string) (error, *models.TokenDetails) {
	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(as.config.Auth.AuthRefreshTokenSalt), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		return errors.New("Refresh token expired"), nil
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err, nil
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return err, nil
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return err, nil
		}
		//Delete the previous Refresh Token
		delErr := as.repo.DeleteToken(refreshUuid)
		if delErr != nil { //if any goes wrong
			return errors.New("unauthorized"), nil
		}

		return as.createToken(userId)
	} else {
		return errors.New("refresh expired"), nil
	}
}

func (as *AuthService) CheckAuthorization(r *http.Request) (error, uint64) {
	tokenAuth, err := as.extractTokenMetadata(r)
	if err != nil {
		return err, 0
	}
	err, userId := as.repo.GetUserIdFromToken(tokenAuth.AccessUuid)
	if err != nil {
		return err, 0
	}
	return nil, userId
}

func (as *AuthService) createToken(userId uint64) (error, *models.TokenDetails) {
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(as.config.Auth.AuthAccessTokenExpire).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(as.config.Auth.AuthRefreshTokenExpire).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(as.config.Auth.AuthAccessTokenSalt))
	if err != nil {
		return err, nil
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(as.config.Auth.AuthRefreshTokenSalt))
	if err != nil {
		return err, nil
	}

	err = as.repo.SaveToken(userId, td)
	if err != nil {
		return err, nil
	}

	return nil, td
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

func (as *AuthService) extractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := as.verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func (as *AuthService) verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := as.extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(as.config.Auth.AuthAccessTokenSalt), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (as *AuthService) extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
