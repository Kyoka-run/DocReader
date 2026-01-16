package main

import (
	"DocReader/internal/handler"
	"DocReader/internal/middleware"
	"DocReader/pkg/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	if err := config.Load("config/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())

	// Routes
	api := r.Group("/api")
	{
		api.POST("/chat", handler.Chat)
		api.POST("/chat/stream", handler.ChatStream)
		api.POST("/upload", handler.FileUpload)
		api.GET("/health", handler.Health)
	}

	port := config.Get().Server.Port
	if port == "" {
		port = "6872"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
