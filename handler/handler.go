package handler

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/scalable-websocket/config"
)

// ScalableWebSocketHandler defines the interface for a websocket handler
type ScalableWebSocketHandler interface {
	HealthCheck(c *fiber.Ctx) error
	StaticFile(c *fiber.Ctx) error
	WebSocket(c *fiber.Ctx) error
}

// scalableWebSocketHandler implements the ScalableWebSocketHandler interface
type scalableWebSocketHandler struct {
	cfg *config.Config
	rdb *redis.Client
}

// NewScalableWebSocketHandler initializes a new handler instance
func NewScalableWebSocketHandler(cfg *config.Config, rdb *redis.Client) ScalableWebSocketHandler {
	return &scalableWebSocketHandler{
		cfg: cfg,
		rdb: rdb,
	}
}

// HealthCheck provides a simple health check endpoint
func (h *scalableWebSocketHandler) HealthCheck(c *fiber.Ctx) error {
	response := map[string]string{
		"ID":      h.cfg.App.ID,
		"Message": "OK",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// StaticFile serves a static HTML file for the WebSocket client
func (h *scalableWebSocketHandler) StaticFile(c *fiber.Ctx) error {
	return c.SendFile("static/websocketClient.html", true)
}

func (h *scalableWebSocketHandler) WebSocket(c *fiber.Ctx) error {
	return websocket.New(func(ws *websocket.Conn) {
		userID := ws.Query("userID")
		if userID == "" {
			log.Println("userID is required but not provided.")
			ws.Close()
			return
		}

		channelName := fmt.Sprintf("chat_messages:%s", userID)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		sub := h.rdb.Subscribe(ctx, channelName)
		defer sub.Close()

		// This goroutine listens for messages from Redis and sends them to the WebSocket client
		go func() {
			ch := sub.Channel()
			for msg := range ch {
				if err := ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
					log.Printf("Error sending message to WebSocket: %v", err)
					return // Exit the goroutine if we can't send to the WebSocket
				}
			}
		}()

		// This for loop reads messages from the WebSocket and publishes them to Redis
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				return // Exit the function if we can't read from the WebSocket
			}

			fmt.Printf("%s sent: %s\n from %s", ws.RemoteAddr(), msg, os.Getenv("APP_ID"))

			// Publish the received message to the Redis channel
			if err := h.rdb.Publish(ctx, channelName, msg).Err(); err != nil {
				log.Printf("Error publishing message to Redis: %v", err)
				return // Exit the function if we can't publish to Redis
			}
		}
	})(c)
}
