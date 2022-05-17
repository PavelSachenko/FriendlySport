package handlers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "test_project/docs"
	"test_project/internal/services"
)

type Handler struct {
	*gin.Engine
	auth services.Auth
	user services.User
	role services.Role
}

func InitHandlers(auth services.Auth, user services.User, role services.Role) *Handler {
	return &Handler{
		Engine: gin.New(),
		auth:   auth,
		user:   user,
		role:   role,
	}
}

func (h *Handler) InitApi() *gin.Engine {
	api := h.Group("api")
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	h.initAuth(api)
	h.initRole(api)
	h.initUser(api)

	return h.Engine
}
