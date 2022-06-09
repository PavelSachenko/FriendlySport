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
	routes := r.Group("/workout").Use(authMiddleware.AuthRequired)
	routes.POST("/create", svc.workoutCreate)

}

func (svc *ServiceClient) workoutCreate(ctx *gin.Context) {
	handlers.WorkoutCreate(ctx, svc.Workout)
}
