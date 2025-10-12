package handlers

import (
	"taskman-backend/internal/middleware"
	"taskman-backend/internal/models"
	"taskman-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ProjectHandler handles project-related requests
type ProjectHandler struct {
	projectService *services.ProjectService
	orgService     *services.OrganizationService
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectService *services.ProjectService, orgService *services.OrganizationService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		orgService:     orgService,
	}
}

// CreateProject handles creating a new project
func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgIDStr := c.Params("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	// Check if user is member
	isMember, _, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember {
		return c.Status(403).JSON(fiber.Map{"error": "Not a member of this organization"})
	}

	var req models.ProjectCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	project, err := h.projectService.CreateProject(&req, orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create project"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Project created successfully",
		"project": project.ToResponse(),
	})
}

// GetProjects handles getting all projects for an organization
func (h *ProjectHandler) GetProjects(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgIDStr := c.Params("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	// Check if user is member
	isMember, _, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember {
		return c.Status(403).JSON(fiber.Map{"error": "Not a member of this organization"})
	}

	projects, err := h.projectService.GetProjectsByOrg(orgID, &userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get projects"})
	}

	return c.JSON(fiber.Map{
		"projects": projects,
	})
}

// GetProject handles getting a specific project
func (h *ProjectHandler) GetProject(c *fiber.Ctx) error {
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

	project, err := h.projectService.GetProjectByID(projectID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Project not found"})
	}

	// Verify project belongs to organization
	if project.OrgID != orgID {
		return c.Status(403).JSON(fiber.Map{"error": "Project does not belong to this organization"})
	}

	// Get assignees
	assignees, err := h.projectService.GetProjectAssignees(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get project assignees"})
	}

	response := project.ToResponse()
	response.Assignees = assignees

	return c.JSON(fiber.Map{
		"project": response,
	})
}

// UpdateProject handles updating a project
func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
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

	// Check if user is project assignee or creator
	project, err := h.projectService.GetProjectByID(projectID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Project not found"})
	}

	if project.OrgID != orgID {
		return c.Status(403).JSON(fiber.Map{"error": "Project does not belong to this organization"})
	}

	isAssignee, err := h.projectService.IsProjectAssignee(projectID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check project assignment"})
	}

	if !isAssignee && project.CreatedBy != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Not assigned to this project"})
	}

	var req models.ProjectUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	updatedProject, err := h.projectService.UpdateProject(projectID, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update project"})
	}

	// Get assignees
	assignees, err := h.projectService.GetProjectAssignees(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get project assignees"})
	}

	response := updatedProject.ToResponse()
	response.Assignees = assignees

	return c.JSON(fiber.Map{
		"message": "Project updated successfully",
		"project": response,
	})
}

// DeleteProject handles deleting a project
func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
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

	// Check if user is project creator or admin
	project, err := h.projectService.GetProjectByID(projectID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Project not found"})
	}

	if project.OrgID != orgID {
		return c.Status(403).JSON(fiber.Map{"error": "Project does not belong to this organization"})
	}

	// Check if user is creator or admin
	isAdmin, _, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if project.CreatedBy != userID && !isAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Not authorized to delete this project"})
	}

	err = h.projectService.DeleteProject(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete project"})
	}

	return c.JSON(fiber.Map{
		"message": "Project deleted successfully",
	})
}

// MoveProject handles moving a project to a different status
func (h *ProjectHandler) MoveProject(c *fiber.Ctx) error {
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

	// Get project to check ownership
	project, err := h.projectService.GetProjectByID(projectID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Project not found"})
	}

	if project.OrgID != orgID {
		return c.Status(403).JSON(fiber.Map{"error": "Project does not belong to this organization"})
	}

	// Check if user is project assignee or creator
	isAssignee, err := h.projectService.IsProjectAssignee(projectID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check project assignment"})
	}

	if !isAssignee && project.CreatedBy != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Not assigned to this project"})
	}

	var req struct {
		Status models.ProjectStatus `json:"status" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.projectService.MoveProject(projectID, req.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to move project"})
	}

	return c.JSON(fiber.Map{
		"message": "Project moved successfully",
	})
}

// BulkMoveProjects handles moving multiple projects to a different status
func (h *ProjectHandler) BulkMoveProjects(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgIDStr := c.Params("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	// Check if user is member
	isMember, _, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember {
		return c.Status(403).JSON(fiber.Map{"error": "Not a member of this organization"})
	}

	var req struct {
		ProjectIDs []uuid.UUID          `json:"project_ids" validate:"required,min=1"`
		Status     models.ProjectStatus `json:"status" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.projectService.BulkMoveProjects(req.ProjectIDs, req.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to move projects"})
	}

	return c.JSON(fiber.Map{
		"message": "Projects moved successfully",
	})
}
