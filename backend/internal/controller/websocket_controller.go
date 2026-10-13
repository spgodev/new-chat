package controller

import (
	"New_Chat/backend/internal/story"
	"New_Chat/backend/internal/ws"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	hub      *ws.Hub
	upgrader websocket.Upgrader
	typing   *story.TypingService
	presence *story.PresenceService
}

func NewWebSocketController(hub *ws.Hub, typing *story.TypingService, presence *story.PresenceService) *WebSocketController {
	return &WebSocketController{
		hub: hub,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		typing:   typing,
		presence: presence,
	}
}

func (c *WebSocketController) Connect(ctx *gin.Context) {
	userID, ok := getCurrentUserID(ctx)
	if !ok {
		return
	}

	conn, err := c.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	client := ws.NewClient(conn)

	if err := client.ConfigureHeartbeat(60 * time.Second); err != nil {
		_ = client.Close()
		return
	}

	c.hub.Add(userID, client)

	slog.Info(
		"websocket connected",
		"user_id", userID,
		"ip", ctx.ClientIP(),
	)

	defer func() { //cleanup в одном defer, чтобы жизненный цикл offline отработал полностью
		c.hub.Remove(userID, client)

		slog.Info(
			"websocket disconnected",
			"user_id", userID,
			"ip", ctx.ClientIP(),
		)

		cleanupCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := c.presence.UserOffline(cleanupCtx, userID); err != nil {
			slog.Error(
				"failed to process websocket disconnect",
				"user_id", userID,
				"error", err,
			)
		}

		if err := client.Close(); err != nil {
			slog.Debug(
				"failed to close websocket client",
				"user_id", userID,
				"error", err,
			)
		}
	}()

	if err := c.presence.UserOnline(ctx.Request.Context(), userID); err != nil { //становление user онлайн
		slog.Error(
			"failed to process websocket connect",
			"user_id", userID,
			"error", err,
		)
	}

	go client.StartPing(25 * time.Second)

	for {
		data, err := client.ReadMessage()
		if err != nil {
			slog.Debug(
				"websocket read stopped",
				"user_id", userID,
				"error", err,
			)
			break
		}

		var event ws.IncomingEvent

		if err = json.Unmarshal(data, &event); err != nil {
			continue
		}

		switch event.Type {
		case "typing":
			var typingData ws.TypingData

			if err := json.Unmarshal(event.Data, &typingData); err != nil {
				continue
			}

			err := c.typing.HandleTyping(
				ctx.Request.Context(),
				userID,
				typingData.ChatID,
				typingData.IsTyping,
			)
			if err != nil {
				continue
			}
		}

	}
}
