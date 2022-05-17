package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initRole(rg *gin.RouterGroup) {
	auth := rg.Group("roles")
	{
		auth.GET("/", h.getRoles)
	}
}

// GetRoles godoc
// @Summary      Get Roles
// @Description  Get all roles system
// @Tags         roles
// @Success      200 {object} models.Role
// @Failure      500
// @Router       /api/roles [get]
func (h *Handler) getRoles(ctx *gin.Context) {
	err, roles := h.role.All()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, roles)
}
