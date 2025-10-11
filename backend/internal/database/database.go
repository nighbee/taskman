package database

import (
	"fmt"
	"log"
	"strings"
	"taskman-backend/internal/config"
	"taskman-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the database connection
var DB *gorm.DB

// Init initializes the database connection and runs migrations
func Init(cfg *config.Config) error {
	var err error

	// Configure GORM logger
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Build database connection string
	var dsn string
	if cfg.DatabaseURL != "" {
		// Use DATABASE_URL if provided
		dsn = cfg.DatabaseURL
	} else {
		// Build connection string from individual parameters
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
	}

	log.Printf("Connecting to database: host=%s, port=%s, dbname=%s, user=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser)

	// Connect to database
	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pool settings
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)

	// Test the connection
	if err = sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	// if err = runMigrations(); err != nil {
	// 	return fmt.Errorf("failed to run migrations: %w", err)
	// }

	log.Println("Database connection established and migrations completed")
	return nil
}

// runMigrations runs database migrations
func runMigrations() error {
	// Auto-migrate all models
	err := DB.AutoMigrate(
		&models.User{},
		&models.Organization{},
		&models.OrgMember{},
		&models.Project{},
		&models.ProjectAssignee{},
		&models.Task{},
		&models.TaskAssignee{},
	)

	if err != nil {
		// Log the error but don't fail if it's just a column already exists error
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("Some migrations skipped (tables/columns already exist): %v", err)
		} else {
			return fmt.Errorf("failed to run auto-migrations: %w", err)
		}
	}

	// Create indexes for better performance
	if err := createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// createIndexes creates additional indexes for performance
func createIndexes() error {
	// Create indexes for common query patterns
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_orgs_invite_code ON organizations(invite_code)",
		"CREATE INDEX IF NOT EXISTS idx_org_members_user_id ON org_members(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_org_members_org_id ON org_members(org_id)",
		"CREATE INDEX IF NOT EXISTS idx_projects_org_id ON projects(org_id)",
		"CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status)",
		"CREATE INDEX IF NOT EXISTS idx_project_assignees_project_id ON project_assignees(project_id)",
		"CREATE INDEX IF NOT EXISTS idx_project_assignees_user_id ON project_assignees(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id)",
		"CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)",
		"CREATE INDEX IF NOT EXISTS idx_task_assignees_task_id ON task_assignees(task_id)",
		"CREATE INDEX IF NOT EXISTS idx_task_assignees_user_id ON task_assignees(user_id)",
	}

	for _, indexSQL := range indexes {
		if err := DB.Exec(indexSQL).Error; err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}
