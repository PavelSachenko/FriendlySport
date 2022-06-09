package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pavel/user_service/config"
	"github.com/pavel/user_service/pkg/db"
	"github.com/pavel/user_service/pkg/pb/auth"
	"github.com/pavel/user_service/pkg/pb/role"
	"github.com/pavel/user_service/pkg/pb/user"
	"github.com/pavel/user_service/pkg/repository"
	"github.com/pavel/user_service/pkg/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Servers struct {
	Auth *auth.Server
	Role *role.Server
	User *user.Server
}

func main() {
	log.Printf("Initial user service config\r\n")
	err, cfg := config.InitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("cfg: %v\r\n", cfg)
	log.Printf("Initial user service tcp server\r\n")
	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		log.Fatalf("Tcp server error: %v\r\n", err)
	}

	err, gRPCServers := initGRPCServices(cfg)
	if err != nil {
		log.Fatalf("Not init gRPC server: %v", err)
	}
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, gRPCServers.Auth)
	role.RegisterRoleServiceServer(grpcServer, gRPCServers.Role)
	user.RegisterUserServiceServer(grpcServer, gRPCServers.User)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

func initGRPCServices(cfg *config.Config) (error, *Servers) {
	log.Printf("Init redis\r\n")
	err, redisDB := db.InitRedis(cfg)
	if err != nil {
		log.Fatalf("Not connected to redis\r\n")
		return err, nil
	}
	log.Printf("Init Postgres\r\n")
	err, postgres := db.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("Not connected to postgres: %v\r\n", err.Error())
		return err, nil
	}

	log.Printf("Init Auth service\r\n")
	authRepo := repository.InitAuthRedis(redisDB, postgres)
	authService := service.NewAuthService(authRepo, cfg)

	log.Printf("Init User service\r\n")
	userRepo := repository.InitUserPostgres(postgres)
	userService := service.InitUserService(userRepo, cfg)

	log.Printf("Init Role service\r\n")
	roleRepo := repository.InitRolePostgres(postgres)
	roleService := service.InitRoleService(roleRepo)

	authServer := auth.InitGRPCUserServer(
		*cfg,
		authService,
		userService,
		roleService,
	)
	roleServer := role.InitGRPCRoleServer(
		roleService,
	)
	userServer := user.InitGRPCUserServer(
		userService,
	)
	return nil, &Servers{
		Auth: authServer,
		Role: roleServer,
		User: userServer,
	}
}
