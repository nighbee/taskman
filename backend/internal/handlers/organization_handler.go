package handlers

import (
	"taskman-backend/internal/middleware"
	"taskman-backend/internal/models"
	"taskman-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// OrganizationHandler handles organization-related requests
type OrganizationHandler struct {
	orgService *services.OrganizationService
}

// NewOrganizationHandler creates a new organization handler
func NewOrganizationHandler(orgService *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

// CreateOrganization handles creating a new organization
func (h *OrganizationHandler) CreateOrganization(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	var req models.OrganizationCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	org, err := h.orgService.CreateOrganization(&req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create organization"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":      "Organization created successfully",
		"organization": org.ToResponse(),
	})
}

// JoinOrganization handles joining an organization via invite code
func (h *OrganizationHandler) JoinOrganization(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	var req models.OrganizationJoinRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Get organization by invite code
	org, err := h.orgService.GetOrganizationByInviteCode(req.InviteCode)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid or expired invite code"})
	}

	// Add user as member
	err = h.orgService.AddMember(org.ID, userID, models.RoleMember)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to join organization"})
	}

	return c.JSON(fiber.Map{
		"message":      "Successfully joined organization",
		"organization": org.ToResponse(),
	})
}

// GetUserOrganizations handles getting user's organizations
func (h *OrganizationHandler) GetUserOrganizations(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgs, err := h.orgService.GetUserOrganizations(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get organizations"})
	}

	return c.JSON(fiber.Map{
		"organizations": orgs,
	})
}

// GetOrganization handles getting a specific organization
func (h *OrganizationHandler) GetOrganization(c *fiber.Ctx) error {
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
	isMember, role, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember {
		return c.Status(403).JSON(fiber.Map{"error": "Not a member of this organization"})
	}

	org, err := h.orgService.GetOrganizationByID(orgID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Organization not found"})
	}

	response := org.ToResponse()
	response.Role = role

	return c.JSON(fiber.Map{
		"organization": response,
	})
}

// GetOrganizationMembers handles getting organization members
func (h *OrganizationHandler) GetOrganizationMembers(c *fiber.Ctx) error {
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

	members, err := h.orgService.GetOrganizationMembers(orgID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get organization members"})
	}

	return c.JSON(fiber.Map{
		"members": members,
	})
}

// RemoveMember handles removing a member from organization (admin only)
func (h *OrganizationHandler) RemoveMember(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgIDStr := c.Params("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	memberIDStr := c.Params("memberId")
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid member ID"})
	}

	// Check if user is admin
	isMember, role, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember || role != models.RoleAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin access required"})
	}

	// Remove member
	err = h.orgService.RemoveMember(orgID, memberID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to remove member"})
	}

	return c.JSON(fiber.Map{
		"message": "Member removed successfully",
	})
}

// UpdateOrganization handles updating organization details
func (h *OrganizationHandler) UpdateOrganization(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	orgIDStr := c.Params("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid organization ID"})
	}

	// Check if user is admin
	isMember, role, err := h.orgService.IsMember(orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check membership"})
	}

	if !isMember || role != models.RoleAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin access required"})
	}

	var req struct {
		Name string `json:"name" validate:"omitempty,min=2,max=100"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Build update map
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}

	if len(updates) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No fields to update"})
	}

	// Update organization
	err = h.orgService.UpdateOrganization(orgID, updates)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update organization"})
	}

	// Get updated organization
	org, err := h.orgService.GetOrganizationByID(orgID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get updated organization"})
	}

	return c.JSON(fiber.Map{
		"message":      "Organization updated successfully",
		"organization": org.ToResponse(),
	})
}
