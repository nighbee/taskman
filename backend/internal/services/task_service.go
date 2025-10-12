package services

import (
	"fmt"
	"taskman-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaskService handles task-related operations
type TaskService struct {
	db *gorm.DB
}

// NewTaskService creates a new task service
func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{db: db}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(req *models.TaskCreateRequest, projectID, createdBy uuid.UUID) (*models.Task, error) {
	task := &models.Task{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Status:      models.TaskStatusNotStarted,
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

	// Insert task
	if err := tx.Create(task).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Add assignees
	for _, assigneeID := range req.AssigneeIDs {
		assignee := &models.TaskAssignee{
			TaskID:     task.ID,
			UserID:     assigneeID,
			AssignedAt: time.Now(),
		}
		if err := tx.Create(assignee).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to add task assignee: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return task, nil
}

// GetTaskByID retrieves a task by ID
func (s *TaskService) GetTaskByID(id uuid.UUID) (*models.Task, error) {
	var task models.Task
	err := s.db.Where("id = ?", id).First(&task).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// GetTasksByProject retrieves all tasks for a project
func (s *TaskService) GetTasksByProject(projectID uuid.UUID, userID *uuid.UUID) ([]models.TaskResponse, error) {
	var queryResults []models.TaskQueryResult

	query := s.db.Table("tasks t").
		Select("t.id, t.project_id, t.name, t.description, t.status, t.created_by, t.deadline, t.created_at, t.updated_at").
		Where("t.project_id = ?", projectID)

	// Order by creation date
	query = query.Order("t.created_at DESC")

	err := query.Scan(&queryResults).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	// Convert to TaskResponse and get assignees for each task
	tasks := make([]models.TaskResponse, len(queryResults))
	for i, result := range queryResults {
		assignees, err := s.GetTaskAssignees(result.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get task assignees: %w", err)
		}

		tasks[i] = models.TaskResponse{
			ID:          result.ID,
			ProjectID:   result.ProjectID,
			Name:        result.Name,
			Description: result.Description,
			Status:      result.Status,
			CreatedBy:   result.CreatedBy,
			Deadline:    result.Deadline,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
			Assignees:   assignees,
		}
	}

	return tasks, nil
}

// UpdateTask updates a task
func (s *TaskService) UpdateTask(id uuid.UUID, req *models.TaskUpdateRequest) (*models.Task, error) {
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

	// Update task
	var task models.Task
	if err := tx.Model(&task).Where("id = ?", id).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Update assignees if provided
	if req.AssigneeIDs != nil {
		// Remove existing assignees
		if err := tx.Where("task_id = ?", id).Delete(&models.TaskAssignee{}).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to remove existing assignees: %w", err)
		}

		// Add new assignees
		for _, assigneeID := range req.AssigneeIDs {
			assignee := &models.TaskAssignee{
				TaskID:     id,
				UserID:     assigneeID,
				AssignedAt: time.Now(),
			}
			if err := tx.Create(assignee).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to add task assignee: %w", err)
			}
		}
	}

	// Get updated task
	if err := tx.Where("id = ?", id).First(&task).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get updated task: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &task, nil
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(id uuid.UUID) error {
	result := s.db.Where("id = ?", id).Delete(&models.Task{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete task: %w", result.Error)
	}

	return nil
}

// GetTaskAssignees retrieves assignees for a task
func (s *TaskService) GetTaskAssignees(taskID uuid.UUID) ([]models.UserResponse, error) {
	var assignees []models.UserResponse

	err := s.db.Table("users u").
		Select("u.id, u.email, u.full_name, u.created_at").
		Joins("JOIN task_assignees ta ON u.id = ta.user_id").
		Where("ta.task_id = ?", taskID).
		Order("ta.assigned_at ASC").
		Scan(&assignees).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get task assignees: %w", err)
	}

	return assignees, nil
}

// IsTaskAssignee checks if a user is assigned to a task
func (s *TaskService) IsTaskAssignee(taskID, userID uuid.UUID) (bool, error) {
	var count int64
	err := s.db.Model(&models.TaskAssignee{}).
		Where("task_id = ? AND user_id = ?", taskID, userID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check task assignment: %w", err)
	}

	return count > 0, nil
}

// MoveTask moves a task to a different status
func (s *TaskService) MoveTask(id uuid.UUID, status models.TaskStatus) error {
	result := s.db.Model(&models.Task{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to move task: %w", result.Error)
	}

	return nil
}

// BulkMoveTasks moves multiple tasks to a different status
func (s *TaskService) BulkMoveTasks(taskIDs []uuid.UUID, status models.TaskStatus) error {
	result := s.db.Model(&models.Task{}).Where("id IN ?", taskIDs).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to bulk move tasks: %w", result.Error)
	}

	return nil
}
