package main

import (
	_ "github.com/lib/pq"
	"github.com/pavel/user_service/config"
	"github.com/pavel/user_service/pkg/db"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/pb"
	"github.com/pavel/user_service/pkg/pb/auth"
	"github.com/pavel/user_service/pkg/pb/role"
	"github.com/pavel/user_service/pkg/pb/user"
	"github.com/pavel/user_service/pkg/repository"
	"github.com/pavel/user_service/pkg/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

type gRPCServices struct {
	Auth *pb.GRPCAuthService
	Role *pb.GRPCRoleService
	User *pb.GRPCUserService
}

func main() {
	logger := logger.GetLogger()
	logger.Info("Initial user service config")
	err, cfg := config.InitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	logger.Infof("cfg: %v", cfg)
	logger.Info("Initial user service tcp server")
	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		logger.Fatalf("Tcp server error: %v", err)
	}

	err, gRPCServers := initGRPCServices(cfg, logger)
	if err != nil {
		logger.Fatalf("Not init gRPC server: %v", err)
	}
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, gRPCServers.Auth)
	role.RegisterRoleServiceServer(grpcServer, gRPCServers.Role)
	user.RegisterUserServiceServer(grpcServer, gRPCServers.User)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}

func initGRPCServices(cfg *config.Config, logger *logger.Logger) (error, *gRPCServices) {
	logger.Info("Init redis")
	err, redisDB := db.InitRedis(cfg)
	if err != nil {
		logger.Fatalf("Not connected to redis: err %v", err)
		return err, nil
	}
	logger.Info("Init Postgres\r\n")
	err, postgres := db.InitPostgres(cfg)
	if err != nil {
		logger.Fatalf("Not connected to postgres: %v", err.Error())
		return err, nil
	}

	logger.Info("Init Auth service")
	authRepo := repository.InitAuthRedis(redisDB, postgres, logger)
	authService := service.NewAuthService(authRepo, cfg)

	logger.Info("Init User service")
	userRepo := repository.InitUserPostgres(postgres, logger)
	userService := service.InitUserService(userRepo, cfg)

	logger.Info("Init Role service")
	roleRepo := repository.InitRolePostgres(postgres, logger)
	roleService := service.InitRoleService(roleRepo)

	authServer := pb.InitGRPCAuthService(
		*cfg,
		authService,
		userService,
		roleService,
		logger,
	)
	roleServer := pb.InitGRPCRoleServer(
		roleService,
		logger,
	)
	userServer := pb.InitGRPCUserService(
		userService,
		logger,
	)
	return nil, &gRPCServices{
		Auth: authServer,
		Role: roleServer,
		User: userServer,
	}
}
