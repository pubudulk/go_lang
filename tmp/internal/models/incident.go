package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IncidentSeverity represents the severity levels for incidents
type IncidentSeverity string

const (
	Low      IncidentSeverity = "low"
	Medium   IncidentSeverity = "medium"
	High     IncidentSeverity = "high"
	Critical IncidentSeverity = "critical"
)

// IncidentStatus represents the status of an incident
type IncidentStatus string

const (
	Open       IncidentStatus = "open"
	InProgress IncidentStatus = "in_progress"
	Resolved   IncidentStatus = "resolved"
	Closed     IncidentStatus = "closed"
)

type NoteType string

const (
	Update        NoteType = "update"
	Investigation NoteType = "investigation"
	Resolution    NoteType = "resolution"
	Communication NoteType = "communication"
)

// Incident represents an incident in the command platform
type Incident struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	IncidentKey int                `json:"incident_key" bson:"incident_key"`
	Title       string             `json:"title" bson:"title" validate:"required,min=3,max=255"`
	Severity    IncidentSeverity   `json:"severity" bson:"severity" validate:"required,oneof=low medium high critical"`
	Status      IncidentStatus     `json:"status" bson:"status" validate:"required,oneof=open in_progress resolved closed"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	Notes       []Note             `json:"notes" bson:"notes"`
	WatchList   []Watcher          `json:"watchlist" bson:"watchlist"`
	CreatedBy   string             `json:"created_by" bson:"created_by"` // Email of the creator
	Description string             `json:"description" bson:"description"`
	Assignee    string             `json:"assignee" bson:"assignee"`
}

// Note represents a note added to an incident
type Note struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Content     string             `json:"content" bson:"content" validate:"required,min=1,max=1000"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	AuthorEmail string             `json:"author_email" bson:"author_email"` // Email of the author
	Type        NoteType           `json:"type" bson:"type" validate:"required,oneof=update investigation resolution communication"`
}

type Watcher struct {
	Email string `json:"email" bson:"email" validate:"required,email"`
}

// CreateIncidentRequest represents the request payload for creating an incident
type CreateIncidentRequest struct {
	Title       string           `json:"title" validate:"required,min=3,max=255"`
	Severity    IncidentSeverity `json:"severity" validate:"required,oneof=low medium high critical"`
	Description string           `json:"description"`
	Notes       []Note           `json:"notes"`
	AuthorEmail string           `json:"author_email" form:"author_email"` // Email of the creator
	Assignee    string           `json:"assignee"`
}

// UpdateIncidentStatusRequest represents the request payload for updating incident status
type UpdateIncidentStatusRequest struct {
	Status      IncidentStatus `json:"status" validate:"required,oneof=open in_progress resolved closed"`
	AuthorEmail string         `json:"author_email" form:"author_email"` // Email of the creator
}

// UpdateIncidentStatusRequest represents the request payload for updating incident status
type UpdateIncidentSeverityRequest struct {
	Severity    IncidentSeverity `json:"severity" validate:"required,oneof=low medium high critical"`
	AuthorEmail string           `json:"author_email" form:"author_email"` // Email of the creator
}

// AddNoteRequest represents the request payload for adding a note to an incident
type AddNoteRequest struct {
	Content     string   `json:"content" validate:"required,min=1,max=1000"`
	AuthorEmail string   `json:"author_email" form:"author_email"` // Email of the creator
	Type        NoteType `json:"type" validate:"required,oneof=update investigation resolution communication"`
}

// ValidSeverities returns a slice of valid severity values
func ValidSeverities() []IncidentSeverity {
	return []IncidentSeverity{
		Low,
		Medium,
		High,
		Critical,
	}
}

// ValidStatuses returns a slice of valid status values
func ValidStatuses() []IncidentStatus {
	return []IncidentStatus{
		Open,
		InProgress,
		Resolved,
		Closed,
	}
}

// IsValidSeverity checks if the provided severity is valid
func (s IncidentSeverity) IsValid() bool {
	for _, severity := range ValidSeverities() {
		if s == severity {
			return true
		}
	}
	return false
}

// IsValidStatus checks if the provided status is valid
func (s IncidentStatus) IsValid() bool {
	for _, status := range ValidStatuses() {
		if s == status {
			return true
		}
	}
	return false
}
