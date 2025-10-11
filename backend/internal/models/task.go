package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusNotStarted TaskStatus = "not-started"
	TaskStatusInProgress TaskStatus = "in-progress"
	TaskStatusDone       TaskStatus = "done"
)

// Task represents a task in the system
type Task struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProjectID   uuid.UUID      `json:"project_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Status      TaskStatus     `json:"status" gorm:"not null;default:'not-started'"`
	CreatedBy   uuid.UUID      `json:"created_by" gorm:"not null"`
	Deadline    *time.Time     `json:"deadline"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Project       Project        `json:"project" gorm:"foreignKey:ProjectID;references:ID"`
	CreatedByUser User           `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:ID"`
	Assignees     []TaskAssignee `json:"assignees" gorm:"foreignKey:TaskID;references:ID"`
}

// TaskAssignee represents a task assignee
type TaskAssignee struct {
	TaskID     uuid.UUID `json:"task_id" gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey"`
	AssignedAt time.Time `json:"assigned_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relationships
	Task Task `json:"task" gorm:"foreignKey:TaskID;references:ID"`
	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

// TaskCreateRequest represents the request to create a new task
type TaskCreateRequest struct {
	Name        string      `json:"name" validate:"required,min=2,max=100"`
	Description string      `json:"description" validate:"max=500"`
	Deadline    *time.Time  `json:"deadline"`
	AssigneeIDs []uuid.UUID `json:"assignee_ids"`
}

// TaskUpdateRequest represents the request to update a task
type TaskUpdateRequest struct {
	Name        *string     `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string     `json:"description,omitempty" validate:"omitempty,max=500"`
	Status      *TaskStatus `json:"status,omitempty"`
	Deadline    *time.Time  `json:"deadline,omitempty"`
	AssigneeIDs []uuid.UUID `json:"assignee_ids,omitempty"`
}

// TaskMoveRequest represents the request to move a task (drag and drop)
type TaskMoveRequest struct {
	Status TaskStatus `json:"status" validate:"required"`
}

// TaskBulkMoveRequest represents the request to move multiple tasks
type TaskBulkMoveRequest struct {
	TaskIDs []uuid.UUID `json:"task_ids" validate:"required,min=1"`
	Status  TaskStatus  `json:"status" validate:"required"`
}

// TaskResponse represents the task data returned to the client
type TaskResponse struct {
	ID          uuid.UUID      `json:"id"`
	ProjectID   uuid.UUID      `json:"project_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      TaskStatus     `json:"status"`
	CreatedBy   uuid.UUID      `json:"created_by"`
	Deadline    *time.Time     `json:"deadline"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Assignees   []UserResponse `json:"assignees"`
}

// ToResponse converts a Task to TaskResponse
func (t *Task) ToResponse() TaskResponse {
	return TaskResponse{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		Name:        t.Name,
		Description: t.Description,
		Status:      t.Status,
		CreatedBy:   t.CreatedBy,
		Deadline:    t.Deadline,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
