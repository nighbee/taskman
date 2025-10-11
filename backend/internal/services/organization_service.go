package services

import (
	"fmt"
	"math/rand"
	"taskman-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrganizationService handles organization-related operations
type OrganizationService struct {
	db *gorm.DB
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(db *gorm.DB) *OrganizationService {
	return &OrganizationService{db: db}
}

// generateInviteCode generates a random 6-character invite code
func generateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// CreateOrganization creates a new organization
func (s *OrganizationService) CreateOrganization(req *models.OrganizationCreateRequest, createdBy uuid.UUID) (*models.Organization, error) {
	// Generate unique invite code
	var inviteCode string
	for {
		inviteCode = generateInviteCode()
		var exists bool
		err := s.db.Model(&models.Organization{}).Select("1").Where("invite_code = ?", inviteCode).First(&exists).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to check invite code: %w", err)
		}
		if err == gorm.ErrRecordNotFound {
			break
		}
	}

	org := &models.Organization{
		Name:          req.Name,
		CreatedBy:     createdBy,
		InviteCode:    inviteCode,
		CodeExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create organization
	if err := tx.Create(org).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Add creator as admin member
	member := &models.OrgMember{
		OrgID:    org.ID,
		UserID:   createdBy,
		Role:     models.RoleAdmin,
		JoinedAt: time.Now(),
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to add creator as admin: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return org, nil
}

// GetOrganizationByID retrieves an organization by ID
func (s *OrganizationService) GetOrganizationByID(id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	err := s.db.Where("id = ?", id).First(&org).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return &org, nil
}

// GetOrganizationByInviteCode retrieves an organization by invite code
func (s *OrganizationService) GetOrganizationByInviteCode(inviteCode string) (*models.Organization, error) {
	var org models.Organization
	err := s.db.Where("invite_code = ? AND code_expires_at > ?", inviteCode, time.Now()).First(&org).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid or expired invite code")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return &org, nil
}

// AddMember adds a user to an organization
func (s *OrganizationService) AddMember(orgID, userID uuid.UUID, role string) error {
	member := &models.OrgMember{
		OrgID:    orgID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}

	// Use ON CONFLICT DO UPDATE for upsert behavior
	err := s.db.Model(&models.OrgMember{}).
		Where("org_id = ? AND user_id = ?", orgID, userID).
		Assign(map[string]interface{}{"role": role}).
		FirstOrCreate(member).Error

	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	return nil
}

// RemoveMember removes a user from an organization
func (s *OrganizationService) RemoveMember(orgID, userID uuid.UUID) error {
	result := s.db.Where("org_id = ? AND user_id = ?", orgID, userID).Delete(&models.OrgMember{})
	if result.Error != nil {
		return fmt.Errorf("failed to remove member: %w", result.Error)
	}

	return nil
}

// GetUserOrganizations retrieves organizations for a user
func (s *OrganizationService) GetUserOrganizations(userID uuid.UUID) ([]models.OrganizationResponse, error) {
	var orgs []models.OrganizationResponse

	err := s.db.Table("organizations o").
		Select("o.id, o.name, o.created_by, o.invite_code, o.code_expires_at, o.created_at, COUNT(om.user_id) as member_count, om.role").
		Joins("LEFT JOIN org_members om ON o.id = om.org_id").
		Where("om.user_id = ?", userID).
		Group("o.id, o.name, o.created_by, o.invite_code, o.code_expires_at, o.created_at, om.role").
		Order("o.created_at DESC").
		Scan(&orgs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user organizations: %w", err)
	}

	return orgs, nil
}

// IsMember checks if a user is a member of an organization
func (s *OrganizationService) IsMember(orgID, userID uuid.UUID) (bool, string, error) {
	var member models.OrgMember
	err := s.db.Where("org_id = ? AND user_id = ?", orgID, userID).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to check membership: %w", err)
	}

	return true, member.Role, nil
}

// GetOrganizationMembers retrieves all members of an organization
func (s *OrganizationService) GetOrganizationMembers(orgID uuid.UUID) ([]models.UserResponse, error) {
	var members []models.UserResponse

	err := s.db.Table("users u").
		Select("u.id, u.email, u.full_name, u.created_at").
		Joins("JOIN org_members om ON u.id = om.user_id").
		Where("om.org_id = ?", orgID).
		Order("om.joined_at ASC").
		Scan(&members).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get organization members: %w", err)
	}

	return members, nil
}

// UpdateOrganization updates an organization
func (s *OrganizationService) UpdateOrganization(id uuid.UUID, updates map[string]interface{}) error {
	result := s.db.Model(&models.Organization{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update organization: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("organization not found")
	}

	return nil
}
