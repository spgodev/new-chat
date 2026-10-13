package middleware

import (
	"New_Chat/backend/internal/story"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebSocketAuth(tokenService *story.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		if token == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		id, err := tokenService.ParseAccessToken(token)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user_id", id)
		ctx.Next()
	}
}
