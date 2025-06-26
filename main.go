package main

import (
	"gofamtree/config"
	"gofamtree/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting GoFamTree API server...")
	
	// Initialize database connection
	config.InitDB()
	defer func() {
		if sqlDB, err := config.DB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Register all routes
	routes.RegisterRoutes()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("API Endpoints:")
	log.Printf("  POST /admin/login - Admin login")
	log.Printf("  POST /admin/register - Admin registration")
	log.Printf("  GET|POST /houses - List houses | Create house")
	log.Printf("  GET|PUT|DELETE /houses/{id} - Get|Update|Delete house")
	log.Printf("  GET|POST /persons - List persons | Create person")
	log.Printf("  GET|PUT|DELETE /persons/{id} - Get|Update|Delete person")
	log.Printf("  GET|POST /relations - List relations | Create relation")
	log.Printf("  GET|PUT|DELETE /relations/{id} - Get|Update|Delete relation")
	log.Printf("  GET /family-tree/{house_id} - Get family tree for house")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
