package main

import (
	"context"
	"database/sql"
	"event-api/internal/domain/event"
	"event-api/internal/domain/event/repository"
	"event-api/internal/entrypoint/http"
	"event-api/pkg/database"
	"event-api/pkg/webserver"
	"log"

	_ "event-api/docs"

	_ "github.com/lib/pq"
)

// @title Event Service API
// @version 1.0
// @description A RESTful API service for managing events built with Go and PostgreSQL
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@eventservice.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	ctx := context.Background()
	// Initialize database
	cfg := database.NewDBConfig()
	db, err := database.NewDBClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)

	eventRp := repository.NewRepository(db)
	eventSvc := event.NewService(eventRp)
	handler := http.NewEventHandler(eventSvc)

	routes := http.SetupRoutes(handler)
	webserver.NewServer(ctx, routes)
}
