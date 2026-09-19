package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setRefreshCookie(ctx *gin.Context, token string) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("refresh_token",
		token,
		30*24*60*60,
		"/auth",
		"",
		false,
		true)
}

func clearRefreshCookie(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("refresh_token",
		"",
		-1,
		"/auth",
		"",
		false,
		true)
}
