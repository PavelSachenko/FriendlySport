package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"user_service/config"
	_ "user_service/docs"
	"user_service/internal/services"
)

type Handler struct {
	*gin.Engine
	cfg  *config.Config
	auth services.Auth
	user services.User
	role services.Role
}

func InitHandlers(cfg *config.Config, auth services.Auth, user services.User, role services.Role) *Handler {
	return &Handler{
		Engine: gin.New(),
		auth:   auth,
		user:   user,
		role:   role,
		cfg:    cfg,
	}
}

func (h *Handler) InitApi() *gin.Engine {
	api := h.Group("api")
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	h.initAuth(api)
	h.initRole(api)
	h.initUser(api)
	api.Use(cors.Default())

	return h.Engine
}
