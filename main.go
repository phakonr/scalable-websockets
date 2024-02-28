package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/scalable-websocket/config"
	"github.com/scalable-websocket/handler"
	"github.com/scalable-websocket/middleware"
	"github.com/scalable-websocket/redisClient"
)

func main() {

	// Initialize the config
	cfg := config.LoadConfig()

	// Initialize the redis
	rdb := redisClient.NewRedisClient(&cfg)
	defer rdb.Close()

	// Setup Go Fiber

	app := fiber.New()

	// Initialize the handler.
	handler := handler.NewScalableWebSocketHandler(&cfg, rdb)

	// Register the WebSocket route with upgrade check and handler.
	// This route is for establishing WebSocket connections for real-time communication.
	app.Get("/ws", middleware.WebSocketUpgrade, handler.WebSocket)

	// Serve the static HTML file for the root route
	app.Get("/", handler.StaticFile)

	// Health Check
	app.Get("/health", handler.HealthCheck)

	// Start server

	// Get the appId from the environment variable
	appPORT := fmt.Sprintf(":%s", cfg.App.Port)
	appID := fmt.Sprintf(":%s", cfg.App.ID)

	fmt.Printf("APPID %s is listening on PORT %s\n", appID, appPORT)

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("server is shutting down...")
		_ = app.Shutdown()
	}()

	// Listen to host:port
	log.Printf("server is starting on %v", cfg.App.Port)
	app.Listen(appPORT)
}
