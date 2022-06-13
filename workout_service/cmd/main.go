package main

import (
	"github.com/pavel/workout_service/config"
	"github.com/pavel/workout_service/pkg/db"
	"github.com/pavel/workout_service/pkg/pb/workout"
	"github.com/pavel/workout_service/pkg/repository"
	"github.com/pavel/workout_service/pkg/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Servers struct {
	Workout *workout.Server
}

func main() {
	log.Printf("Initial user service config\r\n")
	err, cfg := config.InitConfig()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%v\r\b", cfg)

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		log.Fatalf("Tcp server error: %v\r\n", err)
	}

	err, gRPCServers := initGRPCServices(cfg)
	if err != nil {
		log.Fatalf("Not init gRPC server: %v", err)
	}
	grpcServer := grpc.NewServer()
	workout.RegisterWorkoutServiceServer(grpcServer, gRPCServers.Workout)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

func initGRPCServices(cfg *config.Config) (error, *Servers) {

	log.Printf("Init Postgres\r\n")
	err, postgres := db.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("Not connected to postgres: %v\r\n", err.Error())
		return err, nil
	}
	log.Printf("Init Postgres\r\n")
	err, elasticClient := db.InitElastic(cfg)
	if err != nil {
		log.Fatalf("Not connected to postgres: %v\r\n", err.Error())
		return err, nil
	}
	log.Printf("Init Workout service\r\n")
	workoutRepo := repository.InitWorkoutRepo(postgres, elasticClient)
	workoutService := service.InitWorkoutService(workoutRepo)

	log.Printf("Init Workout server\r\n")
	workoutServer := workout.InitGRPCWorkoutServer(
		workoutService,
	)
	return nil, &Servers{
		Workout: workoutServer,
	}
}
