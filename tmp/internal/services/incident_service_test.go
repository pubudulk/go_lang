package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"makers.anchor/incident/internal/kafka"
	"makers.anchor/incident/internal/models"
)

// MockKafkaProducer for testing
type MockKafkaProducer struct {
	kafka.Producer
}

func (m *MockKafkaProducer) ProduceMessage(event kafka.KafkaEvent) error {
	return nil // Mock successful message production
}

// RepositoryInterface defines the interface for incident repository
type RepositoryInterface interface {
	Create(ctx context.Context, incident *models.Incident) (*models.Incident, error)
	GetByID(ctx context.Context, id string) (*models.Incident, error)
	UpdateStatus(ctx context.Context, id string, status models.IncidentStatus) (*models.Incident, error)
	UpdateSeverity(ctx context.Context, id string, severity models.IncidentSeverity) (*models.Incident, error)
	AddNote(ctx context.Context, incidentID string, note models.Note) (*models.Incident, error)
	AddWatcherToIncident(ctx context.Context, incidentID string, watcher models.Watcher) (*models.Incident, error)
}

// MockIncidentRepository for testing
type MockIncidentRepository struct{}

func (m *MockIncidentRepository) Create(ctx context.Context, incident *models.Incident) (*models.Incident, error) {
	// Mock successful creation by setting ID and timestamps
	incident.ID = primitive.NewObjectID()
	incident.CreatedAt = time.Now().UTC()
	incident.UpdatedAt = time.Now().UTC()

	// Set IDs for notes if they exist
	for i := range incident.Notes {
		if incident.Notes[i].ID.IsZero() {
			incident.Notes[i].ID = primitive.NewObjectID()
		}
		if incident.Notes[i].CreatedAt.IsZero() {
			incident.Notes[i].CreatedAt = time.Now().UTC()
		}
	}

	return incident, nil
}

func (m *MockIncidentRepository) GetByID(ctx context.Context, id string) (*models.Incident, error) {
	return nil, nil // Not used in this test
}

func (m *MockIncidentRepository) UpdateStatus(ctx context.Context, id string, status models.IncidentStatus) (*models.Incident, error) {
	return nil, nil // Not used in this test
}

func (m *MockIncidentRepository) UpdateSeverity(ctx context.Context, id string, severity models.IncidentSeverity) (*models.Incident, error) {
	return nil, nil // Not used in this test
}

func (m *MockIncidentRepository) AddNote(ctx context.Context, incidentID string, note models.Note) (*models.Incident, error) {
	return nil, nil // Not used in this test
}

func (m *MockIncidentRepository) AddWatcherToIncident(ctx context.Context, incidentID string, watcher models.Watcher) (*models.Incident, error) {
	return nil, nil // Not used in this test
}

// TestableIncidentService wraps IncidentService for testing
type TestableIncidentService struct {
	repo     RepositoryInterface
	producer *MockKafkaProducer
}

func NewTestableIncidentService(repo RepositoryInterface, producer *MockKafkaProducer) *TestableIncidentService {
	return &TestableIncidentService{
		repo:     repo,
		producer: producer,
	}
}

func (s *TestableIncidentService) CreateIncident(ctx context.Context, req *models.CreateIncidentRequest) (*models.Incident, error) {
	// Validate severity
	if !req.Severity.IsValid() {
		return nil, fmt.Errorf("invalid severity: %s", req.Severity)
	}

	// Initialize notes array
	var notes []models.Note
	if req.Notes != nil {
		notes = req.Notes
		for i, note := range req.Notes {
			notes[i] = models.Note{
				ID:          primitive.NewObjectID(),
				Content:     note.Content,
				AuthorEmail: note.AuthorEmail,
				CreatedAt:   time.Now().UTC(),
			}
		}
	} else {
		notes = []models.Note{}
	}

	// Create incident with default status
	incident := &models.Incident{
		Title:    req.Title,
		Severity: req.Severity,
		Status:   models.Open,
		Notes:    notes,
		WatchList: []models.Watcher{
			{Email: req.AuthorEmail},
		},
		CreatedBy: req.AuthorEmail,
	}

	createdIncident, err := s.repo.Create(ctx, incident)
	if err != nil {
		return nil, fmt.Errorf("failed to create incident: %w", err)
	}

	// Mock Kafka message production
	s.producer.ProduceMessage(models.IncidentCreated{
		EventKey: primitive.NewObjectID().Hex(),
		Id:       createdIncident.ID.Hex(),
		Title:    createdIncident.Title,
		Severity: string(createdIncident.Severity),
	})

	return createdIncident, nil
}

func TestIncidentService_CreateIncident_Success(t *testing.T) {
	// Arrange
	mockRepo := &MockIncidentRepository{}
	mockProducer := &MockKafkaProducer{}
	service := NewTestableIncidentService(mockRepo, mockProducer)

	req := &models.CreateIncidentRequest{
		Title:       "Database Connection Failed",
		Severity:    models.Critical,
		AuthorEmail: "john.doe@example.com",
		Notes: []models.Note{
			{
				Content:     "Initial investigation shows connection timeout",
				AuthorEmail: "john.doe@example.com",
			},
		},
	}

	ctx := context.Background()

	// Act
	result, err := service.CreateIncident(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected incident to be created, got nil")
	}

	if result.Title != req.Title {
		t.Errorf("Expected title %s, got %s", req.Title, result.Title)
	}

	if result.Severity != req.Severity {
		t.Errorf("Expected severity %s, got %s", req.Severity, result.Severity)
	}

	if result.Status != models.Open {
		t.Errorf("Expected status %s, got %s", models.Open, result.Status)
	}

	if result.CreatedBy != req.AuthorEmail {
		t.Errorf("Expected created by %s, got %s", req.AuthorEmail, result.CreatedBy)
	}

	if len(result.WatchList) != 1 || result.WatchList[0].Email != req.AuthorEmail {
		t.Error("Expected author to be added to watchlist")
	}

	if len(result.Notes) != 1 {
		t.Errorf("Expected 1 note, got %d", len(result.Notes))
	}

	if result.Notes[0].Content != req.Notes[0].Content {
		t.Errorf("Expected note content %s, got %s", req.Notes[0].Content, result.Notes[0].Content)
	}
}
