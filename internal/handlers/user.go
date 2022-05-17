package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initUser(rg *gin.RouterGroup) {
	auth := rg.Group("user").Use(h.authorized)
	{
		auth.GET("/:id", h.getUser)
	}
}

type UserID struct {
	ID uint64 `uri:"id" binding:"required"`
}

// GetUser godoc
// @Security BearerAuth
// @Summary      Get User
// @Description  Get User by ID
// @Tags         user
// @Param        id path int true "user id"
// @Success      200 {object} models.User
// @Failure      400 {object} map[string]string
// @Failure      401 string example="unauthorized"
// @Failure      500 {object} map[string]string
// @Router       /api/user/{id} [get]
func (h *Handler) getUser(ctx *gin.Context) {
	var userId UserID
	if err := ctx.ShouldBindUri(&userId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, user := h.user.GetUser(userId.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
