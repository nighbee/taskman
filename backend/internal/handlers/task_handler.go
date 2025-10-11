package handlers

import (
	"taskman-backend/internal/middleware"
	"taskman-backend/internal/models"
	"taskman-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TaskHandler handles task-related requests
type TaskHandler struct {
	taskService    *services.TaskService
	projectService *services.ProjectService
	orgService     *services.OrganizationService
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(taskService *services.TaskService, projectService *services.ProjectService, orgService *services.OrganizationService) *TaskHandler {
	return &TaskHandler{
		taskService:    taskService,
		projectService: projectService,
		orgService:     orgService,
	}
}

// CreateTask handles creating a new task
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgIDStr := c.Params("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	projectIDStr := c.Params("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	// Check if user is member
	isMember, _, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember {
		return c.Status(403).JSON(fiber.Map{"error": "Not a member of this organization"})
	}

	// Check if user is project assignee
	isAssignee, err := h.projectService.IsProjectAssignee(projectID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check project assignment"})
	}

	if !isAssignee {
		return c.Status(403).JSON(fiber.Map{"error": "Not assigned to this project"})
	}

	var req models.TaskCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	task, err := h.taskService.CreateTask(&req, projectID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create task"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Task created successfully",
		"task":    task.ToResponse(),
	})
}

// GetTasks handles getting all tasks for a project
func (h *TaskHandler) GetTasks(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgIDStr := c.Params("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	projectIDStr := c.Params("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	// Check if user is member
	isMember, _, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember {
		return c.Status(403).JSON(fiber.Map{"error": "Not a member of this organization"})
	}

	tasks, err := h.taskService.GetTasksByProject(projectID, &userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get tasks"})
	}

	return c.JSON(fiber.Map{
		"tasks": tasks,
	})
}

// GetTask handles getting a specific task
func (h *TaskHandler) GetTask(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgIDStr := c.Params("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	projectIDStr := c.Params("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	taskIDStr := c.Params("taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	// Check if user is member
	isMember, _, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember {
		return c.Status(403).JSON(fiber.Map{"error": "Not a member of this organization"})
	}

	task, err := h.taskService.GetTaskByID(taskID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	// Verify task belongs to project
	if task.ProjectID != projectID {
		return c.Status(403).JSON(fiber.Map{"error": "Task does not belong to this project"})
	}

	// Get assignees
	assignees, err := h.taskService.GetTaskAssignees(taskID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get task assignees"})
	}

	response := task.ToResponse()
	response.Assignees = assignees

	return c.JSON(fiber.Map{
		"task": response,
	})
}

// UpdateTask handles updating a task
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgIDStr := c.Params("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	projectIDStr := c.Params("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	taskIDStr := c.Params("taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	// Check if user is member
	isMember, _, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember {
		return c.Status(403).JSON(fiber.Map{"error": "Not a member of this organization"})
	}

	// Check if user is task assignee or creator
	task, err := h.taskService.GetTaskByID(taskID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	if task.ProjectID != projectID {
		return c.Status(403).JSON(fiber.Map{"error": "Task does not belong to this project"})
	}

	isAssignee, err := h.taskService.IsTaskAssignee(taskID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check task assignment"})
	}

	if !isAssignee && task.CreatedBy != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Not assigned to this task"})
	}

	var req models.TaskUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	updatedTask, err := h.taskService.UpdateTask(taskID, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update task"})
	}

	// Get assignees
	assignees, err := h.taskService.GetTaskAssignees(taskID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get task assignees"})
	}

	response := updatedTask.ToResponse()
	response.Assignees = assignees

	return c.JSON(fiber.Map{
		"message": "Task updated successfully",
		"task":    response,
	})
}

// DeleteTask handles deleting a task
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgIDStr := c.Params("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	projectIDStr := c.Params("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	taskIDStr := c.Params("taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	// Check if user is member
	isMember, _, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember {
		return c.Status(403).JSON(fiber.Map{"error": "Not a member of this organization"})
	}

	// Check if user is task creator or admin
	task, err := h.taskService.GetTaskByID(taskID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	if task.ProjectID != projectID {
		return c.Status(403).JSON(fiber.Map{"error": "Task does not belong to this project"})
	}

	// Check if user is creator or admin
	isAdmin, _, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if task.CreatedBy != userID && !isAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Not authorized to delete this task"})
	}

	err = h.taskService.DeleteTask(taskID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete task"})
	}

	return c.JSON(fiber.Map{
		"message": "Task deleted successfully",
	})
}

// MoveTask handles moving a task to a different status
func (h *TaskHandler) MoveTask(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgIDStr := c.Params("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	projectIDStr := c.Params("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	taskIDStr := c.Params("taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	// Check if user is member
	isMember, _, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember {
		return c.Status(403).JSON(fiber.Map{"error": "Not a member of this organization"})
	}

	// Get task to verify it belongs to the project
	task, err := h.taskService.GetTaskByID(taskID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	if task.ProjectID != projectID {
		return c.Status(403).JSON(fiber.Map{"error": "Task does not belong to this project"})
	}

	// Check if user is task assignee
	isAssignee, err := h.taskService.IsTaskAssignee(taskID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check task assignment"})
	}

	if !isAssignee {
		return c.Status(403).JSON(fiber.Map{"error": "Not assigned to this task"})
	}

	var req models.TaskMoveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.taskService.MoveTask(taskID, req.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to move task"})
	}

	return c.JSON(fiber.Map{
		"message": "Task moved successfully",
	})
}

// BulkMoveTasks handles moving multiple tasks to a different status
func (h *TaskHandler) BulkMoveTasks(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgIDStr := c.Params("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	projectIDStr := c.Params("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	// Check if user is member
	isMember, _, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember {
		return c.Status(403).JSON(fiber.Map{"error": "Not a member of this organization"})
	}

	var req models.TaskBulkMoveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Validate that all tasks belong to the project
	for _, taskID := range req.TaskIDs {
		task, err := h.taskService.GetTaskByID(taskID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
		}
		if task.ProjectID != projectID {
			return c.Status(403).JSON(fiber.Map{"error": "Task does not belong to this project"})
		}
	}

	err = h.taskService.BulkMoveTasks(req.TaskIDs, req.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to move tasks"})
	}

	return c.JSON(fiber.Map{
		"message": "Tasks moved successfully",
	})
}
