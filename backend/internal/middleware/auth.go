package middleware

import (
	"strings"
	"taskman-backend/internal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtManager *auth.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Set user information in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		return c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens but doesn't require them
func OptionalAuthMiddleware(jwtManager *auth.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Next()
		}

		token := tokenParts[1]
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			return c.Next()
		}

		// Set user information in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		return c.Next()
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userID := c.Locals("user_id")
	if userID == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "User ID not found in context")
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID type in context")
	}

	return userIDUUID, nil
}

// GetUserEmailFromContext extracts user email from context
func GetUserEmailFromContext(c *fiber.Ctx) (string, error) {
	userEmail := c.Locals("user_email")
	if userEmail == nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "User email not found in context")
	}

	userEmailStr, ok := userEmail.(string)
	if !ok {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Invalid user email type in context")
	}

	return userEmailStr, nil
}
