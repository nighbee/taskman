package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProjectStatus represents the status of a project
type ProjectStatus string

const (
	ProjectStatusIdea       ProjectStatus = "idea"
	ProjectStatusInProgress ProjectStatus = "in-progress"
	ProjectStatusFinished   ProjectStatus = "finished"
)

// Project represents a project in the system
type Project struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrgID       uuid.UUID      `json:"org_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Status      ProjectStatus  `json:"status" gorm:"not null;default:'idea'"`
	CreatedBy   uuid.UUID      `json:"created_by" gorm:"not null"`
	Deadline    *time.Time     `json:"deadline"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Organization  Organization      `json:"organization" gorm:"foreignKey:OrgID;references:ID"`
	CreatedByUser User              `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:ID"`
	Assignees     []ProjectAssignee `json:"assignees" gorm:"foreignKey:ProjectID;references:ID"`
	Tasks         []Task            `json:"tasks" gorm:"foreignKey:ProjectID;references:ID"`
}

// ProjectAssignee represents a project assignee
type ProjectAssignee struct {
	ProjectID  uuid.UUID `json:"project_id" gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey"`
	AssignedAt time.Time `json:"assigned_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relationships
	Project Project `json:"project" gorm:"foreignKey:ProjectID;references:ID"`
	User    User    `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

// ProjectCreateRequest represents the request to create a new project
type ProjectCreateRequest struct {
	Name        string      `json:"name" validate:"required,min=2,max=100"`
	Description string      `json:"description" validate:"max=500"`
	Deadline    *time.Time  `json:"deadline"`
	AssigneeIDs []uuid.UUID `json:"assignee_ids"`
}

// ProjectUpdateRequest represents the request to update a project
type ProjectUpdateRequest struct {
	Name        *string        `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string        `json:"description,omitempty" validate:"omitempty,max=500"`
	Status      *ProjectStatus `json:"status,omitempty"`
	Deadline    *time.Time     `json:"deadline,omitempty"`
	AssigneeIDs []uuid.UUID    `json:"assignee_ids,omitempty"`
}

// ProjectResponse represents the project data returned to the client
type ProjectResponse struct {
	ID          uuid.UUID      `json:"id"`
	OrgID       uuid.UUID      `json:"org_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      ProjectStatus  `json:"status"`
	CreatedBy   uuid.UUID      `json:"created_by"`
	Deadline    *time.Time     `json:"deadline"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Assignees   []UserResponse `json:"assignees"`
	TaskCount   int            `json:"task_count"`
}

// ToResponse converts a Project to ProjectResponse
func (p *Project) ToResponse() ProjectResponse {
	return ProjectResponse{
		ID:          p.ID,
		OrgID:       p.OrgID,
		Name:        p.Name,
		Description: p.Description,
		Status:      p.Status,
		CreatedBy:   p.CreatedBy,
		Deadline:    p.Deadline,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
