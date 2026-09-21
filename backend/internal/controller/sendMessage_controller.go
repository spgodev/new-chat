package controller

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/story"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SendMessageReq struct {
	Text string `json:"text"`
}

type SendMessageController struct {
	messages *story.SendMessageService
}

func NewSendMessageController(messages *story.SendMessageService) *SendMessageController {
	return &SendMessageController{messages: messages}
}

func (c *SendMessageController) Send(ctx *gin.Context) {
	userID, ok := getCurrentUserID(ctx)
	if !ok {
		return
	}

	chatID, ok := getChatID(ctx)
	if !ok {
		return
	}
	var req SendMessageReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil || req.Text == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
		return
	}
	message, err := c.messages.Send(ctx.Request.Context(), chatID, userID, req.Text)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrChatNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "chat not found"})

		case errors.Is(err, domain.ErrAccessDenied):
			ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})

		case errors.Is(err, domain.ErrEmptyMessage):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": message})
}
