package main

import (
	"log"
	"os"

	"github.com/bayuwibisana/gofamtree/config"
	"github.com/bayuwibisana/gofamtree/handlers"
	"github.com/bayuwibisana/gofamtree/routes"
	"github.com/bayuwibisana/gofamtree/services"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	config.ConnectDatabase()

	// Initialize services
	personService := services.NewPersonService(config.DB)

	// Initialize handlers with services
	personHandler := handlers.NewPersonHandler(personService)

	// Setup routes with handlers
	r := routes.SetupRoutes(personHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
