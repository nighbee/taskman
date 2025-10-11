package main

import (
	"log"
	"strings"
	"taskman-backend/internal/auth"
	"taskman-backend/internal/config"
	"taskman-backend/internal/database"
	"taskman-backend/internal/handlers"
	"taskman-backend/internal/middleware"
	"taskman-backend/internal/services"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	err := database.Init(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.CORSAllowedOrigins, ","),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
		AllowHeaders:     "Origin,Content-Length,Content-Type,Authorization",
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}))

	// Initialize services
	userService := services.NewUserService(database.GetDB())
	orgService := services.NewOrganizationService(database.GetDB())
	projectService := services.NewProjectService(database.GetDB())
	taskService := services.NewTaskService(database.GetDB())

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, 24*time.Hour)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, jwtManager)
	orgHandler := handlers.NewOrganizationHandler(orgService)
	projectHandler := handlers.NewProjectHandler(projectService, orgService)
	taskHandler := handlers.NewTaskHandler(taskService, projectService, orgService)
	wsHandler := handlers.NewWebSocketHandler()

	// API routes
	api := app.Group("/api/v1")

	// Public routes
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/refresh", authHandler.RefreshToken)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	// Auth routes
	protected.Get("/auth/profile", authHandler.GetProfile)
	protected.Put("/auth/profile", authHandler.UpdateProfile)

	// Organization routes
	protected.Post("/organizations", orgHandler.CreateOrganization)
	protected.Post("/organizations/join", orgHandler.JoinOrganization)
	protected.Get("/organizations", orgHandler.GetUserOrganizations)
	protected.Get("/organizations/:orgId", orgHandler.GetOrganization)
	protected.Get("/organizations/:orgId/members", orgHandler.GetOrganizationMembers)
	protected.Delete("/organizations/:orgId/members/:memberId", orgHandler.RemoveMember)
	protected.Put("/organizations/:orgId", orgHandler.UpdateOrganization)

	// Project routes
	protected.Post("/organizations/:orgId/projects", projectHandler.CreateProject)
	protected.Get("/organizations/:orgId/projects", projectHandler.GetProjects)
	protected.Get("/organizations/:orgId/projects/:projectId", projectHandler.GetProject)
	protected.Put("/organizations/:orgId/projects/:projectId", projectHandler.UpdateProject)
	protected.Delete("/organizations/:orgId/projects/:projectId", projectHandler.DeleteProject)
	protected.Patch("/organizations/:orgId/projects/:projectId/move", projectHandler.MoveProject)
	protected.Patch("/organizations/:orgId/projects/move", projectHandler.BulkMoveProjects)

	// Task routes
	protected.Post("/organizations/:orgId/projects/:projectId/tasks", taskHandler.CreateTask)
	protected.Get("/organizations/:orgId/projects/:projectId/tasks", taskHandler.GetTasks)
	protected.Get("/organizations/:orgId/projects/:projectId/tasks/:taskId", taskHandler.GetTask)
	protected.Put("/organizations/:orgId/projects/:projectId/tasks/:taskId", taskHandler.UpdateTask)
	protected.Delete("/organizations/:orgId/projects/:projectId/tasks/:taskId", taskHandler.DeleteTask)
	protected.Patch("/organizations/:orgId/projects/:projectId/tasks/:taskId/move", taskHandler.MoveTask)
	protected.Patch("/organizations/:orgId/projects/:projectId/tasks/move", taskHandler.BulkMoveTasks)

	// WebSocket route
	protected.Get("/ws", wsHandler.HandleWebSocket)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
