package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
	UserID    uuid.UUID       `json:"user_id,omitempty"`
}

// WebSocket message types
const (
	MessageTypeTaskMoved      = "task_moved"
	MessageTypeProjectMoved   = "project_moved"
	MessageTypeTaskCreated    = "task_created"
	MessageTypeTaskUpdated    = "task_updated"
	MessageTypeTaskDeleted    = "task_deleted"
	MessageTypeProjectCreated = "project_created"
	MessageTypeProjectUpdated = "project_updated"
	MessageTypeProjectDeleted = "project_deleted"
	MessageTypeUserJoined     = "user_joined"
	MessageTypeUserLeft       = "user_left"
	MessageTypeError          = "error"
)

// TaskMovedData represents the data for a task moved event
type TaskMovedData struct {
	TaskID    uuid.UUID `json:"task_id"`
	ProjectID uuid.UUID `json:"project_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	UserID    uuid.UUID `json:"user_id"`
}

// ProjectMovedData represents the data for a project moved event
type ProjectMovedData struct {
	ProjectID uuid.UUID `json:"project_id"`
	OrgID     uuid.UUID `json:"org_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	UserID    uuid.UUID `json:"user_id"`
}

// TaskCreatedData represents the data for a task created event
type TaskCreatedData struct {
	Task   TaskResponse `json:"task"`
	UserID uuid.UUID    `json:"user_id"`
}

// TaskUpdatedData represents the data for a task updated event
type TaskUpdatedData struct {
	Task   TaskResponse `json:"task"`
	UserID uuid.UUID    `json:"user_id"`
}

// TaskDeletedData represents the data for a task deleted event
type TaskDeletedData struct {
	TaskID    uuid.UUID `json:"task_id"`
	ProjectID uuid.UUID `json:"project_id"`
	UserID    uuid.UUID `json:"user_id"`
}

// ProjectCreatedData represents the data for a project created event
type ProjectCreatedData struct {
	Project ProjectResponse `json:"project"`
	UserID  uuid.UUID       `json:"user_id"`
}

// ProjectUpdatedData represents the data for a project updated event
type ProjectUpdatedData struct {
	Project ProjectResponse `json:"project"`
	UserID  uuid.UUID       `json:"user_id"`
}

// ProjectDeletedData represents the data for a project deleted event
type ProjectDeletedData struct {
	ProjectID uuid.UUID `json:"project_id"`
	OrgID     uuid.UUID `json:"org_id"`
	UserID    uuid.UUID `json:"user_id"`
}

// UserJoinedData represents the data for a user joined event
type UserJoinedData struct {
	UserID uuid.UUID `json:"user_id"`
	OrgID  uuid.UUID `json:"org_id"`
}

// UserLeftData represents the data for a user left event
type UserLeftData struct {
	UserID uuid.UUID `json:"user_id"`
	OrgID  uuid.UUID `json:"org_id"`
}

// ErrorData represents the data for an error event
type ErrorData struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
