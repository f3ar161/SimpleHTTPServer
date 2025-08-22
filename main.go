package main

import (
	"log"
	"os"
	"taller_challenge/api"
	"taller_challenge/internal"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
		log.Println("Make sure to set DATABASE_URL environment variable")
	}

	// Connect to PostgreSQL database
	app := internal.ConnectionDB()
	defer app.DB.Close()

	// Create events repository
	eventRepo := internal.NewEventRepository(app.DB)

	// Get server port from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Start HTTP server
	api.StartServer(eventRepo, port)
}
