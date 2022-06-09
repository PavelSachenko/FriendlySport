package user

import (
	"context"
	"github.com/pavel/user_service/pkg/pb/role"
	"github.com/pavel/user_service/pkg/service"
	"net/http"
)

type Server struct {
	user service.User
}

func InitGRPCUserServer(user service.User) *Server {
	return &Server{
		user: user,
	}
}

func (s Server) One(ctx context.Context, request *OneRequest) (*OneResponse, error) {
	err, user := s.user.GetUser(request.UserId)
	if err != nil {
		return &OneResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &OneResponse{
		Status: http.StatusOK,
		User: &User{
			Id:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			Description: user.Description,
			Avatar:      user.Avatar,
			CreatedAt:   user.CreatedAt.Unix(),
			UpdatedAt:   user.UpdatedAt.Unix(),
			Role: &role.Role{
				Id:          user.Role.ID,
				Title:       user.Role.Title,
				Description: user.Role.Description,
			},
		},
	}, nil
}

func (s Server) mustEmbedUnimplementedUserServiceServer() {
	//TODO implement me
	panic("implement me")
}
