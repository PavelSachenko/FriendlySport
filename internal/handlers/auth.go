package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test_project/internal/models"
)

func (h *Handler) initAuth(rg *gin.RouterGroup) {
	auth := rg.Group("auth")
	{
		auth.POST("sign-in", h.signIn)
		auth.POST("sign-up", h.signUp)
		auth.POST("refresh", h.Refresh)
		auth.Use(h.authorized).POST("logout", h.logout)
	}
}

func (h *Handler) test(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, map[string]string{
		"test": "hello",
	})
}

type SignIn struct {
	Password string `json:"password"  binding:"required,min=6,max=60"`
	Email    string `json:"email"   binding:"required,email"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// SigIn godoc
// @Summary      Sign In to the system
// @Description  Create token for login to the system
// @Tags         auth
// @Param		 json/input body SignIn true "sign up body"
// @Success      200 {object} tokenResponse
// @Failure      400 {objects} map[string]string
// @Failure      422 {string} string "Invalid json provided"
// @Router       /api/auth/sign-in [post]
func (h *Handler) signIn(ctx *gin.Context) {
	var si SignIn
	err := ctx.ShouldBindJSON(&si)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  "json decoding : " + err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}
	user := &models.User{
		PasswordHash: si.Password,
		Email:        si.Email,
	}

	err, tokens := h.auth.SigIn(user)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

type SignUp struct {
	Email          string `json:"email"   binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6,max=60"`
	RepeatPassword string `json:"repeat_password" binding:"required,min=6,max=60"`
	RoleId         uint64 `json:"role_id" binding:"required"`
}

// SignUp godoc
// @Summary      Sign Up to the system
// @Description  Create new user account on postgres db and create pairs token for it
// @Tags         auth
// @Param		 json/input body SignUp true "sign up body"
// @Success      200 {object} tokenResponse
// @Failure      400 {object} map[string]string
// @Failure      422 {string} string
// @Router       /api/auth/sign-up [post]
func (h *Handler) signUp(ctx *gin.Context) {
	var su SignUp
	err := ctx.ShouldBind(&su)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  "json decoding : " + err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	if su.RepeatPassword != su.RepeatPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  "Incorrect repeat password",
			"status": http.StatusBadRequest,
		})
		return
	}

	user := &models.User{
		PasswordHash: su.Password,
		Email:        su.Email,
	}

	err, tokens := h.auth.SignUp(user, su.RoleId)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// Logout godoc
// @Summary      Logout from the system
// @Description  Delete access token from redis
// @Security BearerAuth
// @Tags         auth
// @Success      204
// @Failure      401 {string} string "unauthorized"
// @Router       /api/auth/logout [post]
func (h *Handler) logout(ctx *gin.Context) {
	err := h.auth.Logout(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh godoc
// @Summary      Refresh auth tokens
// @Description  Delete refresh token from redis and create new
// @Tags         auth
// @Param		 json/input body refreshRequest true "refresh_token body"
// @Success      200 {object} tokenResponse
// @Failure      422 {string} string
// @Router       /api/auth/refresh [post]
func (h *Handler) Refresh(ctx *gin.Context) {

	var r refreshRequest
	err := ctx.ShouldBind(&r)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	err, tokenDetails := h.auth.RefreshToken(r.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, tokenResponse{
		AccessToken:  tokenDetails.AccessToken,
		RefreshToken: tokenDetails.RefreshToken,
	})
}
