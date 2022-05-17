package main

import (
	_ "github.com/lib/pq"
	"log"
	"test_project/config"
	_ "test_project/docs"
	"test_project/internal/handlers"
	"test_project/internal/repositories"
	"test_project/internal/server"
	"test_project/internal/services"
	"test_project/pkg/databases"
)

// @title           Sport Friendly API
// @version         1.0
// @description     just api for frontend.

// @host      localhost:7000
// @BasePath  /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	err, server := run()
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = server.Run()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

//DB_PORT=5432
//DB_HOST=postgres

func run() (error, *server.Server) {
	err, cfg := config.InitConfig()
	if err != nil {
		log.Fatalf("Not created config")
	}
	err, h := initHandlers(cfg)
	if err != nil {
		log.Fatalf("Not created handlers")
	}
	server := server.InitServer(cfg, h.InitApi())

	return err, server
}

func initHandlers(cfg *config.Config) (error, *handlers.Handler) {
	err, redisDB := databases.InitRedis(cfg)
	if err != nil {
		log.Fatalf("Not connected to redis")
		return err, nil
	}
	err, postgres := databases.InitPostgres(cfg)
	if err != nil {
		log.Fatalf(err.Error())
		return err, nil
	}

	authRepo := repositories.InitAuthRedis(redisDB, postgres)
	authService := services.NewAuthService(authRepo, cfg)

	userRepo := repositories.InitUserPostgres(postgres)
	userService := services.InitUserService(userRepo, cfg)

	roleRepo := repositories.InitRolePostgres(postgres)
	roleService := services.InitRoleService(roleRepo)

	return nil, handlers.InitHandlers(
		cfg,
		authService,
		userService,
		roleService,
	)
}
