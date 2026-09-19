package controller

import (
	"New_Chat/backend/internal/domain"
	"New_Chat/backend/internal/story"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterController struct {
	register *story.RegisterService
}

func NewRegisterController(register *story.RegisterService) *RegisterController {
	return &RegisterController{register: register}
}

func (c *RegisterController) Register(ctx *gin.Context) {
	var req RegisterReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	input := story.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	tokenPair, err := c.register.Register(ctx.Request.Context(), input)
	if err != nil {
		switch {

		case errors.Is(err, domain.ErrInvalidUsernameFormat):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})

		case errors.Is(err, domain.ErrInvalidEmailFormat):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})

		case errors.Is(err, domain.ErrShortPassword):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})

		case errors.Is(err, domain.ErrEmailAlreadyUsed):
			ctx.JSON(http.StatusConflict, gin.H{"error": "email already used"})

		case errors.Is(err, domain.ErrUsernameAlreadyExist):
			ctx.JSON(http.StatusConflict, gin.H{"error": "this username is already taken"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	setRefreshCookie(ctx, tokenPair.RefreshToken)
	ctx.JSON(http.StatusCreated, gin.H{"access_token": tokenPair.AccessToken})
}
