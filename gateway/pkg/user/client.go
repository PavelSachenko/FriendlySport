package user

import (
	"github.com/pavel/gateway/config"
	"github.com/pavel/gateway/pkg/user/pb/auth"
	"github.com/pavel/gateway/pkg/user/pb/role"
	"github.com/pavel/gateway/pkg/user/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

type ServiceClient struct {
	Auth auth.AuthServiceClient
	Role role.RoleServiceClient
	User user.UserServiceClient
}

func InitServiceClient(cfg config.Config) ServiceClient {
	log.Printf("Initial user service")
	mt := metadata.New(map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Credentials": "true", "Access-Control-Allow-Methods": "*", "Access-Control-Allow-Headers": "Content-Type, access-control-allow-origin, access-control-allow-headers, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"})
	cc, err := grpc.Dial(cfg.UserServiceUrl, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.Header(&mt)))
	if err != nil {
		log.Fatalf("Can not connect to user service: %v", err)
	}

	return ServiceClient{
		Auth: auth.NewAuthServiceClient(cc),
		Role: role.NewRoleServiceClient(cc),
		User: user.NewUserServiceClient(cc),
	}
}
