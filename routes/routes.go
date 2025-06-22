package routes

import (
	"github.com/bayuwibisana/gofamtree/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(personHandler *handlers.PersonHandler) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "message": "Family Tree API is running"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Family relationships table (avoiding route conflicts)
		api.GET("/relationships-table", personHandler.GetRelationshipsTable)
		
		persons := api.Group("/persons")
		{
			persons.POST("", personHandler.CreatePerson)
			persons.GET("", personHandler.GetAllPersons)
			persons.GET("/search", personHandler.SearchPersons)
			persons.GET("/:id", personHandler.GetPersonByID)
			persons.PUT("/:id", personHandler.UpdatePerson)
			persons.DELETE("/:id", personHandler.DeletePerson)
			
			persons.GET("/:id/family-tree", personHandler.GetFamilyTree)
			
			// Person with all relationships (detailed view)  
			persons.GET("/:id/relationships", personHandler.GetPersonWithAllRelationships)
		}

		// NEW: Enhanced family tree endpoints
		familyTree := api.Group("/family-tree")
		{
			// Relationships table (like ChatGPT conversation)
			familyTree.GET("/relationships-table", personHandler.GetFamilyRelationshipsTable)
			
			// Family tree by house/surname (simple format)
			familyTree.GET("/house/:house", personHandler.GetFamilyTreeByHouseSimple)
		}
	}

	return r
} 