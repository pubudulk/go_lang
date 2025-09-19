package routes

import (
	"github.com/gofiber/fiber/v2"
	"makers.anchor/incident/internal/database"
	"makers.anchor/incident/internal/handlers"
	"makers.anchor/incident/internal/kafka"
	"makers.anchor/incident/internal/repository"
	"makers.anchor/incident/internal/services"
)

func SetupIncidentRoutes(api fiber.Router, db *database.DB, producer *kafka.Producer) {
	// Initialize repository, service and handler
	incidentRepo := repository.NewIncidentRepository(db.Database)
	incidentService := services.NewIncidentService(incidentRepo, producer)
	incidentHandler := handlers.NewIncidentHandler(incidentService)

	// Incident routes
	incidents := api.Group("/incidents")
	incidents.Get("/", incidentHandler.GetAllIncidents)
	incidents.Post("/", incidentHandler.CreateIncident)
	incidents.Get("/:id", incidentHandler.GetIncidentByID)
	incidents.Put("/:id/status", incidentHandler.UpdateIncidentStatus)
	incidents.Put("/:id/severity", incidentHandler.UpdateIncidentSeverity)
	incidents.Post("/:id/notes", incidentHandler.AddNoteToIncident)
	incidents.Post("/:id/watchlist", incidentHandler.AddWatcherToIncident)
}
