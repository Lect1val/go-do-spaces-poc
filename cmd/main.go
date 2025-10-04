package main

import (
	"log"

	"go-do-spaces-poc/config"
	"go-do-spaces-poc/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env (only in local/dev)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, skipping...")
	}

	cfg := config.LoadConfig()

	r := router.SetupRouter(cfg)

	log.Println("🚀 Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
