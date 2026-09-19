package controller

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/story"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefreshController struct {
	token *story.TokenService
}

func NewRefreshController(token *story.TokenService) *RefreshController {
	return &RefreshController{token: token}
}

func (c *RefreshController) Refresh(ctx *gin.Context) {

	refToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		clearRefreshCookie(ctx)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
		return
	}

	tokenPair, err := c.token.Refresh(ctx.Request.Context(), refToken)
	if err != nil {
		if errors.Is(err, domain.ErrAccessDenied) {
			clearRefreshCookie(ctx)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	setRefreshCookie(ctx, tokenPair.RefreshToken)

	ctx.JSON(http.StatusOK, gin.H{"access_token": tokenPair.AccessToken})
}
