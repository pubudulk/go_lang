package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"makers.anchor/incident/internal/kafka"
	"makers.anchor/incident/internal/models"
	"makers.anchor/incident/internal/repository"
)

// IncidentService handles business logic for incidents
type IncidentService struct {
	repo     *repository.IncidentRepository
	producer *kafka.Producer
}

// NewIncidentService creates a new incident service
func NewIncidentService(repo *repository.IncidentRepository, producer *kafka.Producer) *IncidentService {
	return &IncidentService{
		repo:     repo,
		producer: producer,
	}
}

// CreateIncident creates a new incident
func (s *IncidentService) CreateIncident(ctx context.Context, req *models.CreateIncidentRequest) (*models.Incident, error) {
	// Validate severity
	if !req.Severity.IsValid() {
		return nil, fmt.Errorf("invalid severity: %s", req.Severity)
	}

	// Initialize notes array - handle both cases where req.Notes might exist or not
	var notes []models.Note
	if req.Notes != nil {
		notes = req.Notes
		for i, note := range req.Notes {
			notes[i] = models.Note{
				ID:          primitive.NewObjectID(), // Assign new ObjectID to each note
				Content:     note.Content,
				AuthorEmail: note.AuthorEmail,
				CreatedAt:   time.Now().UTC(),
			}
		}
	} else {
		// Otherwise, initialize empty slice
		notes = []models.Note{}
	}

	// Block Explanation
	// 1. If the mail is empty keep watchlist empty
	// 2. If the mail is not empty, validate the format
	// 3. If the format is valid, add to watchlist
	var watcherList []models.Watcher = []models.Watcher{}
	if strings.Trim(req.AuthorEmail, " ") != "" {
		if conditionErr := s.validateEmail(req.AuthorEmail); conditionErr != nil {
			return nil, fmt.Errorf("invalid email: %w", conditionErr)
		}
		watcherList = append(watcherList, models.Watcher{Email: req.AuthorEmail})
	}

	// Get next incident key
	nextKey, err := s.repo.GetNextIncidentKey(ctx)
	if err != nil {
		log.Printf("Error generating incident key: %v", err)
		return nil, fmt.Errorf("failed to generate incident key: %w", err)
	}

	// Create incident with default status
	incident := &models.Incident{
		IncidentKey: nextKey,
		Title:       req.Title,
		Severity:    req.Severity,
		Status:      models.Open,
		Notes:       notes,
		WatchList:   watcherList,
		CreatedBy:   req.AuthorEmail,
		Description: req.Description,
		Assignee:    req.Assignee,
	}

	createdIncident, err := s.repo.Create(ctx, incident)
	if err != nil {
		log.Printf("Error creating incident: %v", err)
		return nil, fmt.Errorf("failed to create incident: %w", err)
	}

	log.Printf("Created new incident: ID=%s, Title=%s, Severity=%s",
		createdIncident.ID.Hex(), createdIncident.Title, createdIncident.Severity)

	s.producer.ProduceMessage(models.IncidentCreated{
		EventKey: primitive.NewObjectID().Hex(),
		Id:       createdIncident.ID.Hex(),
		Title:    createdIncident.Title,
		Severity: string(createdIncident.Severity),
	})

	return createdIncident, nil
}

// GetByID fetches an incident by its ID
func (s *IncidentService) GetByID(ctx context.Context, id string) (*models.Incident, error) {
	incident, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error fetching incident by ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to get incident: %w", err)
	}

	log.Printf("Fetched incident: ID=%s, Title=%s, Status=%s",
		incident.ID.Hex(), incident.Title, incident.Status)

	return incident, nil
}

// GetAllIncidents fetches all incidents
func (s *IncidentService) GetAllIncidents(ctx context.Context) ([]models.Incident, error) {
	incidents, err := s.repo.GetAllIncidents(ctx)
	if err != nil {
		log.Printf("Error fetching incidents: %v", err)
		return nil, fmt.Errorf("failed to get incidents: %w", err)
	}

	log.Printf("Fetched %d incidents", len(incidents))
	return incidents, nil
}

// UpdateIncidentStatus updates the status of an incident
func (s *IncidentService) UpdateIncidentStatus(ctx context.Context, id string, req *models.UpdateIncidentStatusRequest) (*models.Incident, error) {
	// Validate status
	if !req.Status.IsValid() {
		return nil, fmt.Errorf("invalid status: %s", req.Status)
	}

	// Check if incident exists first
	existingIncident, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("incident not found: %w", err)
	}

	// Validate status transition (optional business rule)
	if err := s.validateStatusTransition(existingIncident.Status, req.Status); err != nil {
		return nil, fmt.Errorf("invalid status transition: %w", err)
	}

	updatedIncident, err := s.repo.UpdateStatus(ctx, existingIncident.ID.Hex(), req.Status)
	if err != nil {
		log.Printf("Error updating incident status: %v", err)
		return nil, fmt.Errorf("failed to update incident status: %w", err)
	}

	if strings.Trim(req.AuthorEmail, " ") != "" {
		_, err = s.AddWatcherToIncident(ctx, id, &models.Watcher{Email: req.AuthorEmail})
		if err != nil {
			log.Printf("Error adding watcher to incident: %v", err)
			return nil, fmt.Errorf("status updated but failed to add watcher to incident: %w", err)
		}
	}
	log.Printf("Updated incident status: ID=%s, Status=%s", id, req.Status)

	s.producer.ProduceMessage(models.IncidentStatusUpdated{
		EventKey: primitive.NewObjectID().Hex(),
		Id:       updatedIncident.ID.Hex(),
		Title:    updatedIncident.Title,
		Status:   string(updatedIncident.Status),
	})

	return updatedIncident, nil
}

// UpdateIncidentSeverity updates the severity of an incident
func (s *IncidentService) UpdateIncidentSeverity(ctx context.Context, id string, req *models.UpdateIncidentSeverityRequest) (*models.Incident, error) {
	// Validate
	if !req.Severity.IsValid() {
		return nil, fmt.Errorf("invalid severity: %s", req.Severity)
	}

	// Check if incident exists first
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("incident not found: %w", err)
	}

	updatedIncident, err := s.repo.UpdateSeverity(ctx, id, req.Severity)
	if err != nil {
		log.Printf("Error updating incident severity: %v", err)
		return nil, fmt.Errorf("failed to update incident severity: %w", err)
	}
	if strings.Trim(req.AuthorEmail, " ") != "" {
		_, err = s.AddWatcherToIncident(ctx, id, &models.Watcher{Email: req.AuthorEmail})
		if err != nil {
			log.Printf("Error adding watcher to incident: %v", err)
			return nil, fmt.Errorf("updated incident severity but failed to add watcher to incident: %w", err)
		}
	}
	log.Printf("Updated incident severity: ID=%s, Severity=%s", id, req.Severity)

	s.producer.ProduceMessage(models.IncidentSeverityUpdated{
		EventKey: primitive.NewObjectID().Hex(),
		Id:       updatedIncident.ID.Hex(),
		Title:    updatedIncident.Title,
		Severity: string(updatedIncident.Severity),
	})

	return updatedIncident, nil
}

// AddNoteToIncident adds a note to an incident
func (s *IncidentService) AddNoteToIncident(ctx context.Context, incidentID string, req *models.AddNoteRequest) (*models.Incident, error) {
	// Check if incident exists first
	existingIncident, err := s.repo.GetByID(ctx, incidentID)
	if err != nil {
		return nil, fmt.Errorf("incident not found: %w", err)
	}

	note := models.Note{
		Content:     req.Content,
		AuthorEmail: req.AuthorEmail,
		Type:        req.Type,
	}

	updatedIncident, err := s.repo.AddNote(ctx, existingIncident.ID.Hex(), note)
	if err != nil {
		log.Printf("Error adding note to incident: %v", err)
		return nil, fmt.Errorf("failed to add note to incident: %w", err)
	}

	log.Printf("Added note to incident: ID=%s, Author=%s", incidentID, req.AuthorEmail)

	s.producer.ProduceMessage(models.IncidentNoteAdded{
		EventKey: primitive.NewObjectID().Hex(),
		Id:       updatedIncident.ID.Hex(),
		Title:    updatedIncident.Title,
		Content:  note.Content,
	})

	return updatedIncident, nil
}

// validateStatusTransition validates if a status transition is allowed
func (s *IncidentService) validateStatusTransition(currentStatus, newStatus models.IncidentStatus) error {
	// Define allowed transitions (this is business logic that can be customized)
	allowedTransitions := map[models.IncidentStatus][]models.IncidentStatus{
		models.Open: {
			models.InProgress,
			models.Resolved,
			models.Closed,
		},
		models.InProgress: {
			models.Open,
			models.Resolved,
			models.Closed,
		},
		models.Resolved: {
			models.Open,
			models.InProgress,
			models.Closed,
		},
		models.Closed: {
			models.Open, // Allow reopening closed incidents
		},
	}

	// Check if the transition is allowed
	if allowedStatuses, exists := allowedTransitions[currentStatus]; exists {
		for _, allowed := range allowedStatuses {
			if allowed == newStatus {
				return nil // Transition is allowed
			}
		}
	}

	return fmt.Errorf("cannot transition from %s to %s", currentStatus, newStatus)
}

// validateEmail validates the format of an email address Using external
func (s *IncidentService) validateEmail(email string) error {
	// Format validation only
	if err := checkmail.ValidateFormat(email); err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}

	// Optional: Also check if host exists (requires network call)
	// if err := checkmail.ValidateHost(email); err != nil {
	// 	return fmt.Errorf("invalid email host: %w", err)
	// }

	return nil
}

// adds a watcher to an incident
func (s *IncidentService) AddWatcherToIncident(ctx context.Context, incidentID string, watcher *models.Watcher) (*models.Incident, error) {
	// Check if incident exists first
	_, err := s.repo.GetByID(ctx, incidentID)
	if err != nil {
		return nil, fmt.Errorf("incident not found: %w", err)
	}

	// Block Explanation
	// 1. If the mail is empty keep watchlist empty
	// 2. If the mail is not empty, validate the format
	// 3. If the format is valid, add to watchlist
	if conditionErr := s.validateEmail(watcher.Email); conditionErr != nil {
		return nil, fmt.Errorf("invalid email: %w", conditionErr)
	}

	updatedIncident, err := s.repo.AddWatcherToIncident(ctx, incidentID, *watcher)
	if err != nil {
		log.Printf("Error adding watcher to incident: %v", err)
		return nil, fmt.Errorf("failed to add watcher to incident: %w", err)
	}

	log.Printf("Added watcher to incident: ID=%s, Email=%s", incidentID, watcher.Email)

	return updatedIncident, nil
}
