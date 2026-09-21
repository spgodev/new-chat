package controller

import (
	"New_Chat/backend/internal/story"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SearchUsersController struct {
	searchService *story.SearchUsersService
}

func NewSearchUsersController(searchService *story.SearchUsersService) *SearchUsersController {
	return &SearchUsersController{searchService: searchService}
}

func (s *SearchUsersController) SearchUsers(ctx *gin.Context) {

	pref := ctx.Query("username")

	normPref := strings.TrimSpace(pref)
	if normPref == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}
	users, err := s.searchService.Search(ctx.Request.Context(), normPref)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}
