package controller

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/story"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetChatHistoryController struct {
	history *story.GetChatHistoryService
}

func NewGetChatHistoryController(history *story.GetChatHistoryService) *GetChatHistoryController {
	return &GetChatHistoryController{history: history}
}

func (g *GetChatHistoryController) GetHistory(ctx *gin.Context) {

	userID, ok := getCurrentUserID(ctx)
	if !ok {
		return
	}

	chatID, ok := getChatID(ctx)
	if !ok {
		return
	}

	messages, err := g.history.GetHistory(ctx.Request.Context(), chatID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrChatNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		case errors.Is(err, domain.ErrAccessDenied):
			ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"messages": messages})
}
