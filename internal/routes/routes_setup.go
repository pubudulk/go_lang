package routes

import (
	"github.com/gofiber/fiber/v2"
	"makers.anchor/incident/internal/database"
	"makers.anchor/incident/internal/kafka"
)

func SetupRoutes(app *fiber.App, db *database.DB, producer *kafka.Producer) {
	// API group
	api := app.Group("/api/v1")

	// Health routes
	SetupHealthRoutes(app)

	// Notification routes
	SetupIncidentRoutes(api, db, producer)
}
