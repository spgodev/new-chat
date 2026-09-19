package main

import (
	"New_Chat/backend/internal/controller"
	"New_Chat/backend/internal/middleware"
	"New_Chat/backend/internal/repository"
	"New_Chat/backend/internal/story"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, `postgresql://postgres:postgres@db:5432/chat`)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pool.Close()
	err = pool.Ping(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	usersRepo := repository.NewUsersRepositoryPG(pool)
	refreshRepo := repository.NewRefreshTokenRepoPG(pool)
	directChatsRepo := repository.NewDirectChatPG(pool)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		fmt.Println("JWT_SECRET is required")
		return
	}

	tokenService := story.NewTokenService(refreshRepo, jwtSecret, 15*time.Minute, 720*time.Hour)
	registerService := story.NewRegisterService(usersRepo, tokenService)
	loginService := story.NewLoginService(usersRepo, tokenService)
	meService := story.NewMeService(usersRepo)
	createChatService := story.NewCreateDirectChatService(
		directChatsRepo,
		usersRepo,
	)
	getChatsService := story.NewGetChatsService(directChatsRepo)

	registerController := controller.NewRegisterController(registerService)
	loginController := controller.NewLoginController(loginService)
	refreshController := controller.NewRefreshController(tokenService)
	logoutController := controller.NewLogoutController(story.NewLogoutService(refreshRepo))
	meController := controller.NewMeController(meService)
	createChatController := controller.NewCreateChatController(
		createChatService,
	)
	getChatsController := controller.NewGetChatsController(getChatsService)

	router := gin.Default()
	auth := router.Group("/auth")

	auth.POST("/register", registerController.Register)
	auth.POST("/login", loginController.Login)
	auth.POST("/refresh", refreshController.Refresh)
	auth.POST("/logout", logoutController.Logout)

	api := router.Group("/api")
	api.Use(middleware.Auth(tokenService))
	api.GET("/me", meController.Me)
	api.POST("/chats", createChatController.Create)
	api.GET("/chats", getChatsController.GetChats)

	err = router.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
