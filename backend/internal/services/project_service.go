package services

import (
	"fmt"
	"taskman-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProjectService handles project-related operations
type ProjectService struct {
	db *gorm.DB
}

// NewProjectService creates a new project service
func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(req *models.ProjectCreateRequest, orgID, createdBy uuid.UUID) (*models.Project, error) {
	// Set default status if not provided
	status := models.ProjectStatusIdea
	if req.Status != "" {
		switch req.Status {
		case "idea":
			status = models.ProjectStatusIdea
		case "in-progress":
			status = models.ProjectStatusInProgress
		case "finished":
			status = models.ProjectStatusFinished
		}
	}

	project := &models.Project{
		OrgID:       orgID,
		Name:        req.Name,
		Description: req.Description,
		Status:      status,
		CreatedBy:   createdBy,
		Deadline:    req.Deadline,
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Insert project
	if err := tx.Create(project).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Add assignees (including creator)
	assigneeIDs := append(req.AssigneeIDs, createdBy)

	// Remove duplicates
	uniqueAssignees := make(map[uuid.UUID]bool)
	for _, id := range assigneeIDs {
		uniqueAssignees[id] = true
	}

	for assigneeID := range uniqueAssignees {
		assignee := &models.ProjectAssignee{
			ProjectID:  project.ID,
			UserID:     assigneeID,
			AssignedAt: time.Now(),
		}
		if err := tx.Create(assignee).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to add project assignee: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return project, nil
}

// GetProjectByID retrieves a project by ID
func (s *ProjectService) GetProjectByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := s.db.Where("id = ?", id).First(&project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &project, nil
}

// GetProjectsByOrg retrieves all projects for an organization
func (s *ProjectService) GetProjectsByOrg(orgID uuid.UUID, userID *uuid.UUID) ([]models.ProjectResponse, error) {
	var queryResults []models.ProjectQueryResult

	query := s.db.Table("projects p").
		Select("p.id, p.org_id, p.name, p.description, p.status, p.created_by, p.deadline, p.created_at, p.updated_at, COUNT(t.id) as task_count").
		Joins("LEFT JOIN tasks t ON p.id = t.project_id").
		Where("p.org_id = ?", orgID).
		Group("p.id, p.org_id, p.name, p.description, p.status, p.created_by, p.deadline, p.created_at, p.updated_at")

	// Order by creation date
	query = query.Order("p.created_at DESC")

	err := query.Scan(&queryResults).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	// Convert to ProjectResponse and get assignees for each project
	projects := make([]models.ProjectResponse, len(queryResults))
	for i, result := range queryResults {
		assignees, err := s.GetProjectAssignees(result.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get project assignees: %w", err)
		}

		projects[i] = models.ProjectResponse{
			ID:          result.ID,
			OrgID:       result.OrgID,
			Name:        result.Name,
			Description: result.Description,
			Status:      result.Status,
			CreatedBy:   result.CreatedBy,
			Deadline:    result.Deadline,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
			Assignees:   assignees,
			TaskCount:   result.TaskCount,
		}
	}

	return projects, nil
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(id uuid.UUID, req *models.ProjectUpdateRequest) (*models.Project, error) {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Build update map
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Deadline != nil {
		updates["deadline"] = *req.Deadline
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Update project
	var project models.Project
	if err := tx.Model(&project).Where("id = ?", id).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	// Update assignees if provided
	if req.AssigneeIDs != nil {
		// Remove existing assignees
		if err := tx.Where("project_id = ?", id).Delete(&models.ProjectAssignee{}).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to remove existing assignees: %w", err)
		}

		// Add new assignees
		for _, assigneeID := range req.AssigneeIDs {
			assignee := &models.ProjectAssignee{
				ProjectID:  id,
				UserID:     assigneeID,
				AssignedAt: time.Now(),
			}
			if err := tx.Create(assignee).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to add project assignee: %w", err)
			}
		}
	}

	// Get updated project
	if err := tx.Where("id = ?", id).First(&project).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get updated project: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &project, nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(id uuid.UUID) error {
	result := s.db.Where("id = ?", id).Delete(&models.Project{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete project: %w", result.Error)
	}

	return nil
}

// GetProjectAssignees retrieves assignees for a project
func (s *ProjectService) GetProjectAssignees(projectID uuid.UUID) ([]models.UserResponse, error) {
	var assignees []models.UserResponse

	err := s.db.Table("users u").
		Select("u.id, u.email, u.full_name, u.created_at").
		Joins("JOIN project_assignees pa ON u.id = pa.user_id").
		Where("pa.project_id = ?", projectID).
		Order("pa.assigned_at ASC").
		Scan(&assignees).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get project assignees: %w", err)
	}

	return assignees, nil
}

// IsProjectAssignee checks if a user is assigned to a project
func (s *ProjectService) IsProjectAssignee(projectID, userID uuid.UUID) (bool, error) {
	var count int64
	err := s.db.Model(&models.ProjectAssignee{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check project assignment: %w", err)
	}

	return count > 0, nil
}

// MoveProject moves a project to a different status
func (s *ProjectService) MoveProject(id uuid.UUID, status models.ProjectStatus) error {
	result := s.db.Model(&models.Project{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to move project: %w", result.Error)
	}

	return nil
}

// BulkMoveProjects moves multiple projects to a different status
func (s *ProjectService) BulkMoveProjects(projectIDs []uuid.UUID, status models.ProjectStatus) error {
	result := s.db.Model(&models.Project{}).Where("id IN ?", projectIDs).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to bulk move projects: %w", result.Error)
	}

	return nil
}
