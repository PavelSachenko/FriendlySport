package workout

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/config"
	"github.com/pavel/gateway/pkg/user"
	"github.com/pavel/gateway/pkg/workout/handlers"
)

func RegisterRoute(r *gin.Engine, cfg config.Config, authSvc user.ServiceClient) {
	svc := InitServiceClient(cfg)
	authMiddleware := user.InitAuthMiddleware(authSvc)
	r.GET("api/workouts", svc.workouts).Use(authMiddleware.AuthRequired)
	routes := r.Group("api/workout").Use(authMiddleware.AuthRequired)
	routes.POST("/create", svc.workoutCreate)
	routes.DELETE("/:id", svc.workoutDelete)
	routes.PUT("/:id", svc.workoutUpdate)
	routes.GET("/recommendation-title", svc.workoutRecommendationTitle)
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
