package controller

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/story"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogoutController struct {
	logout *story.LogoutService
}

func NewLogoutController(logout *story.LogoutService) *LogoutController {
	return &LogoutController{logout: logout}
}

func (s *LogoutController) Logout(ctx *gin.Context) {

	refToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		clearRefreshCookie(ctx)
		ctx.Status(http.StatusNoContent)
		return
	}
	err = s.logout.Logout(ctx.Request.Context(), refToken)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			clearRefreshCookie(ctx)
			ctx.Status(http.StatusNoContent)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	clearRefreshCookie(ctx)
	ctx.Status(http.StatusNoContent)
}
