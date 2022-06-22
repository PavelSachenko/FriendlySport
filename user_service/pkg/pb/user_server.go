package pb

import (
	"context"
	"github.com/pavel/user_service/pkg/adapter"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/pb/user"
	"github.com/pavel/user_service/pkg/service"
	"net/http"
)

type GRPCUserService struct {
	user   service.User
	logger *logger.Logger
	user.UserServiceServer
}

func InitGRPCUserService(user service.User, logger *logger.Logger) *GRPCUserService {
	return &GRPCUserService{
		user:   user,
		logger: logger,
	}
}

func (s GRPCUserService) One(ctx context.Context, request *user.OneRequest) (*user.OneResponse, error) {
	err, u := s.user.GetUser(request.UserId)
	if err != nil {
		s.logger.Error(err)
		return &user.OneResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}

	return &user.OneResponse{
		Status: http.StatusOK,
		User:   adapter.UserToGRPC(u, s.logger),
	}, nil
}
