package models

import (
	"encoding/json"
)

const (
	SOURCE_SERVICE = "incident"
	EVENT_TOPIC    = "anchor.incident.events"
)

type IncidentCreated struct {
	EventKey      string `json:"event_key"`
	Id            string `json:"id"`
	Title         string `json:"title"`
	Severity      string `json:"severity"`
	SourceService string `json:"source_service"`
	Version       int    `json:"version"`
	EventType     string `json:"event_type"`
}

type IncidentStatusUpdated struct {
	EventKey      string `json:"event_key"`
	Id            string `json:"id"`
	Title         string `json:"title"`
	Status        string `json:"status"`
	SourceService string `json:"source_service"`
	Version       int    `json:"version"`
	EventType     string `json:"event_type"`
}

type IncidentSeverityUpdated struct {
	EventKey      string `json:"event_key"`
	Id            string `json:"id"`
	Title         string `json:"title"`
	Severity      string `json:"severity"`
	SourceService string `json:"source_service"`
	Version       int    `json:"version"`
	EventType     string `json:"event_type"`
}

type IncidentNoteAdded struct {
	EventKey      string `json:"event_key"`
	Id            string `json:"id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	SourceService string `json:"source_service"`
	Version       int    `json:"version"`
	EventType     string `json:"event_type"`
}

func (e IncidentCreated) GetTopic() string {
	return EVENT_TOPIC
}

func (e IncidentCreated) GetEventType() string {
	return "incident.created"
}

func (e IncidentCreated) GetVersion() int {
	return 1
}

func (e IncidentCreated) GetPayload() ([]byte, error) {
	e.Version = e.GetVersion()
	e.EventType = e.GetEventType()
	e.SourceService = SOURCE_SERVICE
	return json.Marshal(e)
}

// Incident Status Updated
func (e IncidentStatusUpdated) GetTopic() string {
	return EVENT_TOPIC
}

func (e IncidentStatusUpdated) GetEventType() string {
	return "incident.status.updated"
}

func (e IncidentStatusUpdated) GetVersion() int {
	return 1
}

func (e IncidentStatusUpdated) GetPayload() ([]byte, error) {
	e.Version = e.GetVersion()
	e.EventType = e.GetEventType()
	e.SourceService = SOURCE_SERVICE
	return json.Marshal(e)
}

// Incident Status Updated
func (e IncidentSeverityUpdated) GetTopic() string {
	return EVENT_TOPIC
}

func (e IncidentSeverityUpdated) GetEventType() string {
	return "incident.severity.updated"
}

func (e IncidentSeverityUpdated) GetVersion() int {
	return 1
}

func (e IncidentSeverityUpdated) GetPayload() ([]byte, error) {
	e.Version = e.GetVersion()
	e.EventType = e.GetEventType()
	e.SourceService = SOURCE_SERVICE
	return json.Marshal(e)
}

// Incident Note Added
func (e IncidentNoteAdded) GetTopic() string {
	return EVENT_TOPIC
}

func (e IncidentNoteAdded) GetEventType() string {
	return "incident.notes.added"
}

func (e IncidentNoteAdded) GetVersion() int {
	return 1
}

func (e IncidentNoteAdded) GetPayload() ([]byte, error) {
	e.Version = e.GetVersion()
	e.EventType = e.GetEventType()
	e.SourceService = SOURCE_SERVICE
	return json.Marshal(e)
}
