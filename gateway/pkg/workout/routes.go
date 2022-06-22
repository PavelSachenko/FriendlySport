package workout

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/config"
	"github.com/pavel/gateway/pkg/user"
	"github.com/pavel/gateway/pkg/workout/handlers"
)

func RegisterRoute(r *gin.RouterGroup, cfg config.Config, authSvc user.ServiceClient) {
	svc := InitServiceClient(cfg)
	authMiddleware := user.InitAuthMiddleware(authSvc)
	r.GET("workouts", svc.workouts).Use(authMiddleware.AuthRequired)
	workouts := r.Group("workout")
	workouts.Use(authMiddleware.AuthRequired)
	workouts.POST("/create", svc.workoutCreate)
	workouts.DELETE("/:workout_id", svc.workoutDelete)
	workouts.PUT("/:workout_id", svc.workoutUpdate)
	workouts.GET("/recommendation-title", svc.workoutRecommendationTitle)

	exercise := workouts.Group(":workout_id/exercise").Use(authMiddleware.AuthRequired)
	exercise.POST("/", svc.exerciseCreate)
	exercise.DELETE("/:exercise_id", svc.exerciseDelete)
	exercise.PUT("/:exercise_id", svc.exerciseUpdate)
	exercise.GET("/recommendation-title", svc.exerciseRecommendationTitle)
}

func (svc *ServiceClient) workoutCreate(ctx *gin.Context) {
	handlers.WorkoutCreate(ctx, svc.Workout)
}

func (svc *ServiceClient) workoutUpdate(ctx *gin.Context) {
	handlers.WorkoutUpdate(ctx, svc.Workout)
}

func (svc *ServiceClient) workoutDelete(ctx *gin.Context) {
	handlers.WorkoutDelete(ctx, svc.Workout)
}

func (svc *ServiceClient) workouts(ctx *gin.Context) {
	handlers.WorkoutsAll(ctx, svc.Workout)
}

func (svc *ServiceClient) workoutRecommendationTitle(ctx *gin.Context) {
	handlers.WorkoutRecommendationTitle(ctx, svc.Workout)
}

func (svc *ServiceClient) exerciseCreate(ctx *gin.Context) {
	handlers.ExerciseCreate(ctx, svc.Exercise)
}
func (svc *ServiceClient) exerciseUpdate(ctx *gin.Context) {
	handlers.ExerciseUpdate(ctx, svc.Exercise)
}

func (svc *ServiceClient) exerciseDelete(ctx *gin.Context) {
	handlers.ExerciseDelete(ctx, svc.Exercise)
}

func (svc *ServiceClient) exerciseRecommendationTitle(ctx *gin.Context) {
	handlers.ExerciseRecommendationTitle(ctx, svc.Exercise)
}
