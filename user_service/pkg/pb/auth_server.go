package pb

import (
	"context"
	"github.com/pavel/user_service/config"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/model"
	"github.com/pavel/user_service/pkg/pb/auth"
	"github.com/pavel/user_service/pkg/service"
	"net/http"
)

type GRPCAuthService struct {
	auth   service.Auth
	user   service.User
	role   service.Role
	cfg    config.Config
	logger *logger.Logger
	auth.AuthServiceServer
}

func InitGRPCAuthService(
	cfg config.Config,
	auth service.Auth,
	user service.User,
	role service.Role,
	logger *logger.Logger) *GRPCAuthService {
	return &GRPCAuthService{
		cfg:    cfg,
		auth:   auth,
		user:   user,
		role:   role,
		logger: logger,
	}
}

func (a GRPCAuthService) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	err, res := a.auth.SignUp(&model.User{
		Email:        request.Email,
		PasswordHash: request.Password,
	}, request.RoleId)
	if err != nil {
		return &auth.RegisterResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}
	return &auth.RegisterResponse{
		Status:       http.StatusCreated,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (a GRPCAuthService) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	err, res := a.auth.SigIn(&model.User{
		PasswordHash: request.Password,
		Email:        request.Email,
	})
	if err != nil {
		return &auth.LoginResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}

	return &auth.LoginResponse{
		Status:             http.StatusOK,
		AccessToken:        res.AccessToken,
		RefreshToken:       res.RefreshToken,
		RefreshTokenExpire: uint64(a.cfg.Auth.AuthRefreshTokenExpire),
	}, nil
}

func (a GRPCAuthService) CheckAuthToken(ctx context.Context, request *auth.CheckTokenRequest) (*auth.CheckTokenResponse, error) {
	err, userId := a.auth.CheckAuthorization(request.Token)
	if err != nil {
		return &auth.CheckTokenResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}

	return &auth.CheckTokenResponse{
		Status: http.StatusOK,
		UserId: userId,
	}, nil
}

func (a GRPCAuthService) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	err := a.auth.Logout(request.Token)
	if err != nil {
		return &auth.LogoutResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}

	return &auth.LogoutResponse{
		Status: http.StatusNoContent,
	}, nil
}

func (a GRPCAuthService) RefreshAuthToken(ctx context.Context, request *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	err, tokens := a.auth.RefreshToken(request.RefreshToken)
	if err != nil {
		return &auth.RefreshTokenResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}

	return &auth.RefreshTokenResponse{
		Status:             http.StatusOK,
		RefreshToken:       tokens.RefreshToken,
		AccessToken:        tokens.AccessToken,
		RefreshTokenExpire: uint64(a.cfg.Auth.AuthRefreshTokenExpire),
	}, nil
}
