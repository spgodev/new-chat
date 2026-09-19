package middleware

import (
	"New_Chat/backend/internal/story"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(tokenService *story.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		scheme, token, found := strings.Cut(header, " ")
		if !found || scheme != "Bearer" || len(token) == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userId, err := tokenService.ParseAccessToken(token)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("user_id", userId)
		ctx.Next()
		
	}
}
