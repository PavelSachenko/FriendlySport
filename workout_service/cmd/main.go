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

	//test := context.WithValue(context.Background(), "id", "asd")
	//id := uint64(test.Value("id").(int))
	//fmt.Println(id)
	//return
	log.Printf("Initial user service config\r\n")
	err, cfg := config.InitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	//err, postgres := db.InitPostgres(cfg, db.InitPostgresQueryBuilder())
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	////repo := repository.InitExerciseRepo(postgres)
	//ctx := context.WithValue(context.Background(), "user_id", 1)
	//ctx = context.WithValue(ctx, "workout_id", 3)
	//ctx = context.WithValue(ctx, "id", 3)
	//fmt.Println(ctx.Value("user_id"))
	////return
	////temp := "test"
	////err, res := repo.Update(ctx, model.ExerciseUpdate{Description: &temp})
	////if err != nil {
	////	log.Fatalln(err)
	////}
	//repo := repository.InitWorkoutRepo(postgres, nil)
	//err, res := repo.All(1, model.WorkoutsFiltering{Limit: 10, Offset: 0})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(res[0])
	//return
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
	err, postgres := db.InitPostgres(cfg, db.InitPostgresQueryBuilder())
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
