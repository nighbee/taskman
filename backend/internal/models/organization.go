package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organization represents an organization in the system
type Organization struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name          string         `json:"name" gorm:"not null"`
	CreatedBy     uuid.UUID      `json:"created_by" gorm:"not null"`
	InviteCode    string         `json:"invite_code" gorm:"uniqueIndex;not null"`
	CodeExpiresAt time.Time      `json:"code_expires_at" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	CreatedByUser User        `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:ID"`
	Members       []OrgMember `json:"members" gorm:"foreignKey:OrgID;references:ID"`
	Projects      []Project   `json:"projects" gorm:"foreignKey:OrgID;references:ID"`
}

// OrgMember represents a member of an organization
type OrgMember struct {
	OrgID    uuid.UUID `json:"org_id" gorm:"type:uuid;primaryKey"`
	UserID   uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey"`
	Role     string    `json:"role" gorm:"not null;default:'member'"`
	JoinedAt time.Time `json:"joined_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relationships
	Organization Organization `json:"organization" gorm:"foreignKey:OrgID;references:ID"`
	User         User         `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

// OrganizationCreateRequest represents the request to create a new organization
type OrganizationCreateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

// OrganizationJoinRequest represents the request to join an organization
type OrganizationJoinRequest struct {
	InviteCode string `json:"invite_code" validate:"required,len=6"`
}

// OrganizationResponse represents the organization data returned to the client
type OrganizationResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	CreatedBy     uuid.UUID `json:"created_by"`
	InviteCode    string    `json:"invite_code"`
	CodeExpiresAt time.Time `json:"code_expires_at"`
	CreatedAt     time.Time `json:"created_at"`
	MemberCount   int       `json:"member_count"`
	Role          string    `json:"role,omitempty"`
}

// ToResponse converts an Organization to OrganizationResponse
func (o *Organization) ToResponse() OrganizationResponse {
	return OrganizationResponse{
		ID:            o.ID,
		Name:          o.Name,
		CreatedBy:     o.CreatedBy,
		InviteCode:    o.InviteCode,
		CodeExpiresAt: o.CodeExpiresAt,
		CreatedAt:     o.CreatedAt,
	}
}

// MemberRole constants
const (
	RoleMember = "member"
	RoleAdmin  = "admin"
)
