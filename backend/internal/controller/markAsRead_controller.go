package controller

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/story"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MarkAsReadController struct {
	markService *story.MarkAsReadService
}

func NewMarkAsReadController(markService *story.MarkAsReadService) *MarkAsReadController {
	return &MarkAsReadController{markService: markService}
}

func (c *MarkAsReadController) MarkAsRead(ctx *gin.Context) {
	userID, ok := getCurrentUserID(ctx)
	if !ok {
		return
	}

	chatID, ok := getChatID(ctx)
	if !ok {
		return
	}

	err := c.markService.MarkAsRead(ctx.Request.Context(), chatID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrChatNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "chat not found"})

		case errors.Is(err, domain.ErrAccessDenied):
			ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	ctx.Status(http.StatusNoContent)
}
