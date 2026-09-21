package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getCurrentUserID(ctx *gin.Context) (int64, bool) {
	rawUserID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "access denied"})
		return 0, false
	}

	userID, ok := rawUserID.(int64)
	if !ok || userID <= 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't get userID"})
		return 0, false
	}
	return userID, true
}

func getChatID(ctx *gin.Context) (int64, bool) {
	rawChatID := ctx.Param("chat_id")
	chatID, err := strconv.ParseInt(rawChatID, 10, 64)
	if err != nil || chatID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "can't get chatID"})
		return 0, false
	}
	return chatID, true
}
