package role

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pavel/user_service/pkg/service"
	"net/http"
)

type Server struct {
	role service.Role
}

func InitGRPCRoleServer(role service.Role) *Server {
	return &Server{
		role: role,
	}
}
func (s Server) All(ctx context.Context, empty *empty.Empty) (*AllResponse, error) {
	err, roles := s.role.All()
	if err != nil {
		return &AllResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}
	var responseRole []*Role
	for _, role := range roles {
		newRole := &Role{Id: role.ID, Title: role.Title, Description: role.Description}
		responseRole = append(responseRole, newRole)
	}

	return &AllResponse{
		Status: http.StatusOK,
		Roles:  responseRole,
	}, nil
}

func (s Server) mustEmbedUnimplementedRoleServiceServer() {
	//TODO implement me
	panic("implement me")
}
