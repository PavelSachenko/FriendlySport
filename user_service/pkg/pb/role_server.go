package pb

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pavel/user_service/pkg/adapter"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/pb/role"
	"github.com/pavel/user_service/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type GRPCRoleService struct {
	role   service.Role
	logger *logger.Logger
	role.RoleServiceServer
}

func InitGRPCRoleServer(role service.Role, logger *logger.Logger) *GRPCRoleService {
	return &GRPCRoleService{
		role:   role,
		logger: logger,
	}
}
func (s GRPCRoleService) All(ctx context.Context, empty *empty.Empty) (*role.AllResponse, error) {
	mt := metadata.New(map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Credentials": "true", "Access-Control-Allow-Methods": "*", "Access-Control-Allow-Headers": "Content-Type, access-control-allow-origin, access-control-allow-headers, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"})
	err := grpc.SendHeader(ctx, mt)
	if err != nil {
		s.logger.Error(err)
	}
	err, roles := s.role.All()
	if err != nil {
		s.logger.Error(err.Error())
		return &role.AllResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  err.Error(),
		}, nil
	}

	return &role.AllResponse{
		Status: http.StatusOK,
		Roles:  adapter.RoleListToGRPC(roles, s.logger),
	}, nil
}
