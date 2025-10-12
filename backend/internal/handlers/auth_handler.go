package handlers

import (
	"strings"
	"taskman-backend/internal/auth"
	"taskman-backend/internal/middleware"
	"taskman-backend/internal/models"
	"taskman-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	userService *services.UserService
	jwtManager  *auth.JWTManager
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService *services.UserService, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.UserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if user already exists
	_, err := h.userService.GetUserByEmail(req.Email)
	if err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "User already exists"})
	}

	// Create user
	user, err := h.userService.CreateUser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User created successfully",
		"user":    user.ToResponse(),
		"token":   token,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Get user by email
	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Validate password
	if err := h.userService.ValidatePassword(user, req.Password); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user.ToResponse(),
		"token":   token,
	})
}

// GetMe handles getting current user info
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"user": user.ToResponse(),
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Authorization header required"})
	}

	// Extract token from "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization header format"})
	}

	token := tokenParts[1]
	newToken, err := h.jwtManager.RefreshToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	return c.JSON(fiber.Map{
		"message": "Token refreshed successfully",
		"token":   newToken,
	})
}

// GetProfile handles getting user profile
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"user": user.ToResponse(),
	})
}

// UpdateProfile handles updating user profile
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	var req struct {
		FullName string `json:"full_name" validate:"omitempty,min=2,max=100"`
		Email    string `json:"email" validate:"omitempty,email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Build update map
	updates := make(map[string]interface{})
	if req.FullName != "" {
		updates["full_name"] = req.FullName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	if len(updates) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No fields to update"})
	}

	// Update user
	err = h.userService.UpdateUser(userID, updates)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	// Get updated user
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get updated user"})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user":    user.ToResponse(),
	})
}
