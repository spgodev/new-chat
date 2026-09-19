package controller

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/story"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MeController struct {
	meService *story.MeService
}

func NewMeController(meService *story.MeService) *MeController {
	return &MeController{meService: meService}
}

func (m *MeController) Me(ctx *gin.Context) {
	value, exist := ctx.Get("user_id")
	if !exist {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID, ok := value.(int64)
	if !ok || userID <= 0 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	me, err := m.meService.Me(ctx.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"user_id":  me.ID,
		"email":    me.Email,
		"username": me.Username,
	})
}
