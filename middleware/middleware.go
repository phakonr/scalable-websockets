package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// WebSocketUpgrade checks for WebSocket upgrade and upgrades the connection
func WebSocketUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
