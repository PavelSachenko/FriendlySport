package main

import (
	"fmt"
	"github.com/pavel/workout_service/config"
	"github.com/pavel/workout_service/pkg/db"
	"github.com/pavel/workout_service/pkg/logger"
	"github.com/pavel/workout_service/pkg/pb"
	"github.com/pavel/workout_service/pkg/pb/workout"
	"github.com/pavel/workout_service/pkg/repository"
	"github.com/pavel/workout_service/pkg/service"
	"google.golang.org/grpc"
	"net"
)

type gRPCServices struct {
	Workout  *pb.GRPCWorkoutService
	Exercise *pb.GRPCExerciseService
}

func main() {
	logger := logger.GetLogger()
	logger.Info("test")
	cfg := getConfig(logger)
	lis := getTCPServer(logger, cfg)
	gRPCServices := getGRPCServices(logger, cfg)
	InitGrpcServer(logger, lis, gRPCServices)
}

func getConfig(logger *logger.Logger) *config.Config {
	logger.Info("Init config")
	err, cfg := config.InitConfig(logger)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed init config. ERROR: %v", err))
	}
	return cfg
}

func getTCPServer(logger *logger.Logger, cfg *config.Config) net.Listener {
	logger.Info("Init tcp server")
	lis, err := net.Listen("tcp", cfg.Server.Host+cfg.Server.Port)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed init tcp server. ERROR: %s", err))
	}
	logger.Info(fmt.Sprintf("GRPCWorkoutService address: %s", lis.Addr().Network()))
	return lis
}

func getGRPCServices(logger *logger.Logger, cfg *config.Config) *gRPCServices {
	logger.Info(fmt.Sprintf("Init gRPC services"))
	err, gRPCServices := initGRPCServices(logger, cfg)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed init gRPC services. ERROR: %s", err))
	}
	return gRPCServices
}

func initGRPCServices(logger *logger.Logger, cfg *config.Config) (error, *gRPCServices) {

	logger.Info("Init Postgres DB")
	err, postgres := db.InitPostgres(cfg, db.InitPostgresQueryBuilder())
	if err != nil {
		logger.Fatal(fmt.Sprintf("Not connected to postgres. ERROR: %s", err.Error()))
	}
	logger.Info("Init ElasticSearch DB")
	err, elasticClient := db.InitElastic(cfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Not connected to ElasticSearch. ERROR: %s", err.Error()))
		return err, nil
	}

	logger.Info("Init workout service")
	workoutRepo := repository.InitWorkoutPostgresRepo(logger, postgres, elasticClient)
	workoutService := service.InitWorkoutService(logger, workoutRepo)

	logger.Info("Init workout gGRPC service")
	workoutGRPCService := pb.InitGRPCWorkoutService(
		workoutService,
		logger,
	)

	exerciseRepo := repository.InitExercisePostgresRepo(logger, postgres, elasticClient)
	exerciseService := service.InitExerciseService(exerciseRepo)
	logger.Info("Init exercise gGRPC service")
	exerciseGRPCService := pb.InitGRPCExerciseService(exerciseService, logger)
	return nil, &gRPCServices{
		Workout:  workoutGRPCService,
		Exercise: exerciseGRPCService,
	}
}

func InitGrpcServer(logger *logger.Logger, lis net.Listener, services *gRPCServices) {
	logger.Info(fmt.Sprintf("Init gRPC server"))
	grpcServer := grpc.NewServer()
	workout.RegisterWorkoutServiceServer(grpcServer, services.Workout)
	workout.RegisterExerciseServiceServer(grpcServer, services.Exercise)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Info(fmt.Sprintf("Failed serve gRPC server. ERROR: %s", err))
	}
}
