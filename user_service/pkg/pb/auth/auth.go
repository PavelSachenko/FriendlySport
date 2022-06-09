package auth

import (
	"context"
	"github.com/pavel/user_service/config"
	"github.com/pavel/user_service/pkg/model"
	"github.com/pavel/user_service/pkg/service"
	"net/http"
)

type Server struct {
	auth service.Auth
	user service.User
	role service.Role
	cfg  config.Config
}

func InitGRPCUserServer(cfg config.Config, auth service.Auth, user service.User, role service.Role) *Server {
	return &Server{
		cfg:  cfg,
		auth: auth,
		user: user,
		role: role,
	}
}

func (s Server) Register(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	err, res := s.auth.SignUp(&model.User{
		Email:        request.Email,
		PasswordHash: request.Password,
	}, request.RoleId)
	if err != nil {
		return &RegisterResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}
	return &RegisterResponse{
		Status:       http.StatusCreated,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (s Server) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	err, res := s.auth.SigIn(&model.User{
		PasswordHash: request.Password,
		Email:        request.Email,
	})
	if err != nil {
		return &LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &LoginResponse{
		Status:             http.StatusOK,
		AccessToken:        res.AccessToken,
		RefreshToken:       res.RefreshToken,
		RefreshTokenExpire: uint64(s.cfg.Auth.AuthRefreshTokenExpire),
	}, nil
}

func (s Server) CheckAuthToken(ctx context.Context, request *CheckTokenRequest) (*CheckTokenResponse, error) {
	err, userId := s.auth.CheckAuthorization(request.Token)
	if err != nil {
		return &CheckTokenResponse{
			Status: http.StatusInternalServerError,
		}, nil
	}

	return &CheckTokenResponse{
		Status: http.StatusOK,
		UserId: userId,
	}, nil
}

func (s Server) Logout(ctx context.Context, request *LogoutRequest) (*LogoutResponse, error) {
	err := s.auth.Logout(request.Token)
	if err != nil {
		return &LogoutResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &LogoutResponse{
		Status: http.StatusNoContent,
	}, nil
}

func (s Server) RefreshAuthToken(ctx context.Context, request *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	err, tokens := s.auth.RefreshToken(request.RefreshToken)
	if err != nil {
		return &RefreshTokenResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &RefreshTokenResponse{
		Status:             http.StatusOK,
		RefreshToken:       tokens.RefreshToken,
		AccessToken:        tokens.AccessToken,
		RefreshTokenExpire: uint64(s.cfg.Auth.AuthRefreshTokenExpire),
	}, nil
}

func (s Server) mustEmbedUnimplementedAuthServiceServer() {
	//TODO implement me
	panic("implement me")
}
