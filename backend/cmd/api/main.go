package main

import (
	"New_Chat/backend/internal/controller"
	"New_Chat/backend/internal/middleware"
	"New_Chat/backend/internal/repository"
	"New_Chat/backend/internal/story"
	"New_Chat/backend/internal/ws"
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	logger.Info("application starting")

	ctx := context.Background()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		logger.Error("DB_URL is required")
		return
	}
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		return
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		return
	}

	logger.Info("database connected")

	usersRepo := repository.NewUsersRepositoryPG(pool)
	refreshRepo := repository.NewRefreshTokenRepoPG(pool)
	directChatsRepo := repository.NewDirectChatPG(pool)
	messageRepo := repository.NewMessageRepositoryPG(pool)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Error("JWT_SECRET is required")
		return
	}

	hub := ws.NewHub()

	tokenService := story.NewTokenService(refreshRepo, jwtSecret, 15*time.Minute, 720*time.Hour)
	registerService := story.NewRegisterService(usersRepo, tokenService)
	loginService := story.NewLoginService(usersRepo, tokenService)
	meService := story.NewMeService(usersRepo)
	createChatService := story.NewCreateDirectChatService(
		directChatsRepo,
		usersRepo,
		hub,
	)
	logoutService := story.NewLogoutService(refreshRepo)
	getChatsService := story.NewGetChatsService(directChatsRepo, usersRepo, hub)
	searchUsersService := story.NewSearchUsersService(usersRepo)
	sendMessagesService := story.NewSendMessageService(messageRepo, directChatsRepo, hub)
	getChatHistoryService := story.NewGetChatHistoryService(messageRepo, directChatsRepo)
	markAsReadService := story.NewMarkAsReadService(
		messageRepo,
		directChatsRepo,
		hub,
	)
	typingService := story.NewTypingService(directChatsRepo, hub)
	presenceService := story.NewPresenceService(directChatsRepo, usersRepo, hub)

	registerController := controller.NewRegisterController(registerService)
	loginController := controller.NewLoginController(loginService)
	refreshController := controller.NewRefreshController(tokenService)
	logoutController := controller.NewLogoutController(logoutService)
	meController := controller.NewMeController(meService)
	createChatController := controller.NewCreateChatController(
		createChatService,
	)
	getChatsController := controller.NewGetChatsController(getChatsService)
	searchUsersController := controller.NewSearchUsersController(searchUsersService)
	sendMessagesController := controller.NewSendMessageController(sendMessagesService)
	getChatHistoryController := controller.NewGetChatHistoryController(getChatHistoryService)
	markAsReadController := controller.NewMarkAsReadController(
		markAsReadService,
	)
	webSocketController := controller.NewWebSocketController(hub, typingService, presenceService)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5500",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PATCH",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	router.GET(
		"/api/ws",
		middleware.WebSocketAuth(tokenService),
		webSocketController.Connect,
	)

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
	api.GET("/users/search", searchUsersController.SearchUsers)

	chats := api.Group("/chats")
	chats.POST("/:chat_id/messages", sendMessagesController.Send)
	chats.GET("/:chat_id/messages", getChatHistoryController.GetHistory)
	chats.PATCH(
		"/:chat_id/messages/read",
		markAsReadController.MarkAsRead,
	)

	logger.Info("http server starting", "addr", ":8080")
	if err := router.Run(":8080"); err != nil {
		logger.Error("http server stopped", "error", err)
	}
}
