package user

import (
	"github.com/gin-gonic/gin"
	"github.com/pavel/gateway/config"
	"github.com/pavel/gateway/pkg/user/handlers"
)

func RegisterRoute(r *gin.RouterGroup, cfg config.Config) ServiceClient {
	svc := InitServiceClient(cfg)
	authMiddleware := AuthMiddleware{svc: svc}
	r.GET("roles", svc.roles)
	r.GET("testt", func(context *gin.Context) {
		context.JSON(200, "test")
	})
	routes := r.Group("/auth")
	routes.POST("login", svc.login)
	routes.POST("register", svc.register)
	routes.POST("logout", svc.logout).Use(authMiddleware.AuthRequired)
	routes.POST("refresh", svc.refresh).Use(authMiddleware.AuthRequired)
	r.Use(authMiddleware.AuthRequired).GET("user", svc.user)

	return svc
}

func (svc *ServiceClient) roles(ctx *gin.Context) {
	handlers.Roles(ctx, svc.Role)
}
func (svc *ServiceClient) login(ctx *gin.Context) {
	handlers.Login(ctx, svc.Auth)
}
func (svc *ServiceClient) register(ctx *gin.Context) {
	handlers.Register(ctx, svc.Auth)
}
func (svc *ServiceClient) logout(ctx *gin.Context) {
	handlers.Logout(ctx, svc.Auth)
}
func (svc *ServiceClient) refresh(ctx *gin.Context) {
	handlers.RefreshAuthToken(ctx, svc.Auth)
}
func (svc *ServiceClient) user(ctx *gin.Context) {
	handlers.User(ctx, svc.User)
}
