package controller

import (
	"New_Chat/backend/internal/story"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetChatsController struct {
	chats *story.GetChatsService
}

func NewGetChatsController(chats *story.GetChatsService) *GetChatsController {
	return &GetChatsController{chats: chats}
}

func (c *GetChatsController) GetChats(ctx *gin.Context) {

	currentUserID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "something wrong"})
		return
	}

	userID, ok := currentUserID.(int64)
	if !ok || userID <= 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "something wrong"})
		return
	}
	chats, err := c.chats.GetByUser(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"chats": chats})
}
