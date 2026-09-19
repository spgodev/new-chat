package controller

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/story"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginController struct {
	login *story.LoginService
}

func NewLoginController(login *story.LoginService) *LoginController {
	return &LoginController{login: login}
}

func (c *LoginController) Login(ctx *gin.Context) {
	var req LoginReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	input := story.LoginInput{
		Username: req.Username,
		Password: req.Password,
	}
	tokenPair, err := c.login.Login(ctx.Request.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	setRefreshCookie(ctx, tokenPair.RefreshToken)
	ctx.JSON(http.StatusOK, gin.H{
		"access_token": tokenPair.AccessToken,
	})
}
