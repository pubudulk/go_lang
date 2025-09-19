package handlers

import (
	"github.com/gofiber/fiber/v2"
	"makers.anchor/incident/internal/models"
	"makers.anchor/incident/internal/services"
)

// IncidentHandler handles HTTP requests for incidents
type IncidentHandler struct {
	service *services.IncidentService
}

// NewIncidentHandler creates a new incident handler
func NewIncidentHandler(service *services.IncidentService) *IncidentHandler {
	return &IncidentHandler{
		service: service,
	}
}

// CreateIncident handles POST /incidents
func (h *IncidentHandler) CreateIncident(c *fiber.Ctx) error {
	var req models.CreateIncidentRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Basic validation
	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	if req.Severity == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Severity is required",
		})
	}

	incident, err := h.service.CreateIncident(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create incident",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    incident,
	})
}

// GetAllIncidents handles GET /incidents
func (h *IncidentHandler) GetAllIncidents(c *fiber.Ctx) error {
	incidents, err := h.service.GetAllIncidents(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve incidents",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    incidents,
	})
}

// GetIncidentByID handles GET /incidents/:id
func (h *IncidentHandler) GetIncidentByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incident ID is required",
		})
	}

	incident, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "incident not found" || err.Error() == "no documents found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Incident not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve incident",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    incident,
	})
}

// UpdateIncidentStatus handles PUT /incidents/:id/status
func (h *IncidentHandler) UpdateIncidentStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incident ID is required",
		})
	}

	var req models.UpdateIncidentStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status is required",
		})
	}

	incident, err := h.service.UpdateIncidentStatus(c.Context(), id, &req)
	if err != nil {
		if err.Error() == "incident not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Incident not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to update incident status",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    incident,
	})
}

// UpdateIncidentSeverity handles PUT /incidents/:id/severity
func (h *IncidentHandler) UpdateIncidentSeverity(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incident ID is required",
		})
	}

	var req models.UpdateIncidentSeverityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if req.Severity == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Severity is required",
		})
	}

	incident, err := h.service.UpdateIncidentSeverity(c.Context(), id, &req)
	if err != nil {
		if err.Error() == "incident not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Incident not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to update incident severity",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    incident,
	})
}

// AddNoteToIncident handles POST /incidents/:id/notes
func (h *IncidentHandler) AddNoteToIncident(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incident ID is required",
		})
	}

	var req models.AddNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Note content is required",
		})
	}

	incident, err := h.service.AddNoteToIncident(c.Context(), id, &req)
	if err != nil {
		if err.Error() == "incident not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Incident not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to add note to incident",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    incident,
	})
}

// AddWatcherToIncident
func (h *IncidentHandler) AddWatcherToIncident(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incident ID is required",
		})
	}

	var req models.Watcher
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	incident, err := h.service.AddWatcherToIncident(c.Context(), id, &req)
	if err != nil {
		if err.Error() == "incident not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Incident not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to add watcher to incident",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    incident,
	})
}
