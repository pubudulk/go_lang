package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"makers.anchor/incident/internal/models"
)

const (
	IncidentsCollection = "incidents"
)

// IncidentRepository handles incident database operations
type IncidentRepository struct {
	collection *mongo.Collection
}

// NewIncidentRepository creates a new incident repository
func NewIncidentRepository(db *mongo.Database) *IncidentRepository {
	return &IncidentRepository{
		collection: db.Collection(IncidentsCollection),
	}
}

// Create creates a new incident
func (r *IncidentRepository) Create(ctx context.Context, incident *models.Incident) (*models.Incident, error) {
	// Set timestamps
	now := time.Now()
	incident.CreatedAt = now
	incident.UpdatedAt = now
	incident.ID = primitive.NewObjectID()

	// Initialize empty notes slice if nil
	if incident.Notes == nil {
		incident.Notes = []models.Note{}
	}

	// Insert the incident
	result, err := r.collection.InsertOne(ctx, incident)
	if err != nil {
		return nil, fmt.Errorf("failed to create incident: %w", err)
	}

	// Set the generated ID
	incident.ID = result.InsertedID.(primitive.ObjectID)

	return incident, nil
}

// UpdateStatus updates the status of an incident
func (r *IncidentRepository) UpdateStatus(ctx context.Context, id string, status models.IncidentStatus) (*models.Incident, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid incident ID format: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedIncident models.Incident
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, opts).Decode(&updatedIncident)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("incident not found")
		}
		return nil, fmt.Errorf("failed to update incident status: %w", err)
	}

	return &updatedIncident, nil
}

// UpdateSeverity updates the severity of an incident
func (r *IncidentRepository) UpdateSeverity(ctx context.Context, id string, severity models.IncidentSeverity) (*models.Incident, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid incident ID format: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"severity":   severity,
			"updated_at": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedIncident models.Incident
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, opts).Decode(&updatedIncident)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("incident not found")
		}
		return nil, fmt.Errorf("failed to update incident severity: %w", err)
	}

	return &updatedIncident, nil
}

// AddNote adds a note to an incident
func (r *IncidentRepository) AddNote(ctx context.Context, incidentID string, note models.Note) (*models.Incident, error) {
	objectID, err := primitive.ObjectIDFromHex(incidentID)
	if err != nil {
		return nil, fmt.Errorf("invalid incident ID format: %w", err)
	}

	// Set note metadata
	note.ID = primitive.NewObjectID()
	note.CreatedAt = time.Now()

	update := bson.M{
		"$push": bson.M{"notes": note},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedIncident models.Incident
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, opts).Decode(&updatedIncident)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("incident not found")
		}
		return nil, fmt.Errorf("failed to add note to incident: %w", err)
	}

	return &updatedIncident, nil
}

func (r *IncidentRepository) GetByID(ctx context.Context, id string) (*models.Incident, error) {
	// Convert string ID to integer
	incidentKey, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid incident ID format: %w", err)
	}

	var incident models.Incident
	err = r.collection.FindOne(ctx, bson.M{"incident_key": incidentKey}).Decode(&incident)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("incident not found")
		}
		return nil, fmt.Errorf("failed to get incident: %w", err)
	}

	return &incident, nil
}

// GetAll retrieves all incidents with optional filtering and pagination
func (r *IncidentRepository) GetAllIncidents(ctx context.Context) ([]models.Incident, error) {
	opts := options.Find()

	// Sort by created_at descending (newest first)
	opts.SetSort(bson.D{bson.E{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get incidents: %w", err)
	}
	defer cursor.Close(ctx)

	var incidents []models.Incident
	err = cursor.All(ctx, &incidents)
	if err != nil {
		return nil, fmt.Errorf("failed to decode incidents: %w", err)
	}

	return incidents, nil
}

// Add add watcher to an incident
func (r *IncidentRepository) AddWatcherToIncident(ctx context.Context, incidentID string, watcher models.Watcher) (*models.Incident, error) {
	objectID, err := primitive.ObjectIDFromHex(incidentID)
	if err != nil {
		return nil, fmt.Errorf("invalid incident ID: %w", err)
	}

	update := bson.M{
		"$addToSet": bson.M{"watchlist": watcher},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedIncident models.Incident
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, opts).Decode(&updatedIncident)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("incident not found")
		}
		return nil, fmt.Errorf("failed to add watcher to incident: %w", err)
	}

	return &updatedIncident, nil
}

// GetNextIncidentKey gets the next auto-increment ID for incidents
func (r *IncidentRepository) GetNextIncidentKey(ctx context.Context) (int, error) {
	// Find the incident with the highest IncidentKey
	opts := options.FindOne().SetSort(bson.D{bson.E{Key: "incident_key", Value: -1}})

	var incident models.Incident
	err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&incident)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No incidents exist, start from 1
			return 1, nil
		}
		return 0, fmt.Errorf("failed to get max incident key: %w", err)
	}

	// Return the next ID
	return incident.IncidentKey + 1, nil
}
