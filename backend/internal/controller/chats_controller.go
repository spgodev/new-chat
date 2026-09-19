package controller

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/story"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CreateReq struct {
	OtherUsername string `json:"username"`
}

type CreateChatController struct {
	createChats *story.CreateDirectChatService
}

func NewCreateChatController(createChats *story.CreateDirectChatService) *CreateChatController {
	return &CreateChatController{createChats: createChats}
}

func (c *CreateChatController) Create(ctx *gin.Context) {
	var req CreateReq

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

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	req.OtherUsername = strings.TrimSpace(req.OtherUsername)
	if req.OtherUsername == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	input := story.CreateDirectChatInput{
		CurrentUserID: userID,
		OtherUsername: req.OtherUsername,
	}
	chat, err := c.createChats.Create(ctx.Request.Context(), input)
	if err != nil {

		switch {

		case errors.Is(err, domain.ErrUserNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})

		case errors.Is(err, domain.ErrCantCreateChat):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot be created with oneself"})

		case errors.Is(err, domain.ErrChatAlreadyExist):
			ctx.JSON(http.StatusConflict, gin.H{"error": "chat already exists"})

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"chat": chat})
}
