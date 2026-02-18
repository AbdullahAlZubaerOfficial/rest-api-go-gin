package main

import (
	"database/sql"
	"fmt"
	"log"

	"rest-api-in-gin/internal/database"
	"rest-api-in-gin/internal/env"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	models := database.NewModels(db)

	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}

func (app *application) serve() error {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.Default()

	// Register routes
	router.GET("/health", app.healthCheck)
	
	// Event routes
	router.GET("/events", app.getAllEvents)
	router.GET("/events/:id", app.getEvent)
	router.POST("/events", app.createEvent)
	router.PUT("/events/:id", app.updateEvent)
	router.DELETE("/events/:id", app.deleteEvent)

	// Start server
	log.Printf("Starting server on port %d", app.port)
	return router.Run(fmt.Sprintf(":%d", app.port))
}

// Health check handler
func (app *application) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

// Event handlers
func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, events)
}

func (app *application) getEvent(c *gin.Context) {
	// Parse ID from URL
	var id int
	_, err := fmt.Sscanf(c.Param("id"), "%d", &id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if event == nil {
		c.JSON(404, gin.H{"error": "Event not found"})
		return
	}
	c.JSON(200, event)
}

func (app *application) createEvent(c *gin.Context) {
	var event database.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := app.models.Events.Insert(&event); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, event)
}

func (app *application) updateEvent(c *gin.Context) {
	// Parse ID from URL
	var id int
	_, err := fmt.Sscanf(c.Param("id"), "%d", &id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	var event database.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	// Set the ID from URL
	event.Id = id

	if err := app.models.Events.Update(&event); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, event)
}

func (app *application) deleteEvent(c *gin.Context) {
	// Parse ID from URL
	var id int
	_, err := fmt.Sscanf(c.Param("id"), "%d", &id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Event deleted successfully"})
}