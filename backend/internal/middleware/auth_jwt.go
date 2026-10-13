package middleware

import (
	"New_Chat/backend/internal/story"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(tokenService *story.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		header := ctx.GetHeader("Authorization")
		scheme, token, found := strings.Cut(header, " ")

		if !found || scheme != "Bearer" || len(token) == 0 {
			slog.Warn("authorization header is invalid",
				"method", ctx.Request.Method,
				"path", ctx.Request.URL.Path,
				"ip", ctx.ClientIP(),
			)

			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err := tokenService.ParseAccessToken(token)
		if err != nil {
			slog.Warn(
				"access token validation failed",
				"method", ctx.Request.Method,
				"path", ctx.Request.URL.Path,
				"ip", ctx.ClientIP(),
				"error", err,
			)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("user_id", userID)
		ctx.Next()

	}
}
