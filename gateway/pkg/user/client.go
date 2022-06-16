package user

import (
	"github.com/pavel/gateway/config"
	"github.com/pavel/gateway/pkg/user/pb/auth"
	"github.com/pavel/gateway/pkg/user/pb/role"
	"github.com/pavel/gateway/pkg/user/pb/user"
	"google.golang.org/grpc"
	"log"
)

type ServiceClient struct {
	Auth auth.AuthServiceClient
	Role role.RoleServiceClient
	User user.UserServiceClient
}

func InitServiceClient(cfg config.Config) ServiceClient {
	log.Printf("Initial user service")
	cc, err := grpc.Dial(cfg.UserServiceUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can not connect to user service: %v", err)
	}
	return ServiceClient{
		Auth: auth.NewAuthServiceClient(cc),
		Role: role.NewRoleServiceClient(cc),
		User: user.NewUserServiceClient(cc),
	}
}
