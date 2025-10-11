# TaskMan Backend - Fiber & GORM Version

A collaborative task management backend built with Go, Fiber, and GORM.

## Features

- **User Authentication**: JWT-based authentication with registration and login
- **Organization Management**: Create organizations with invite codes, manage members
- **Project Management**: Kanban-style project boards with drag-and-drop support
- **Task Management**: Nested tasks within projects with assignee management
- **Real-time Collaboration**: WebSocket support for live updates
- **Role-based Permissions**: Admin and member roles with appropriate access control
- **Auto-migrations**: GORM handles database schema automatically

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Fiber (high-performance HTTP framework)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens
- **Real-time**: WebSocket connections
- **Containerization**: Docker

## Project Structure

```
backend/
├── internal/
│   ├── auth/           # JWT authentication
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and auto-migrations
│   ├── handlers/       # HTTP request handlers (Fiber)
│   ├── middleware/     # Custom middleware (Fiber)
│   ├── models/         # Data models with GORM tags
│   └── services/       # Business logic layer with GORM
├── migrations/         # SQL migration files (reference)
├── main.go            # Application entry point
├── Dockerfile         # Container configuration
├── go.mod            # Go module dependencies
└── env.example       # Environment variables template
```

## Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Docker (optional)

### Environment Variables

Copy `env.example` to `.env` and configure:

```bash
# Database Configuration
DATABASE_URL=postgresql://username:password@localhost:5432/taskman
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRY=24h

# Server Configuration
PORT=8080
GIN_MODE=debug

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### Database Setup

1. Create a PostgreSQL database
2. The application will automatically run migrations on startup using GORM
3. Update the `DATABASE_URL` in your `.env` file

### Running the Application

#### Development

```bash
# Install dependencies
go mod download

# Run the application
go run main.go
```

#### Docker

```bash
# Build the image
docker build -t taskman-backend .

# Run the container
docker run -p 8080:8080 --env-file .env taskman-backend
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `GET /api/v1/auth/profile` - Get user profile
- `PUT /api/v1/auth/profile` - Update user profile

### Organizations

- `POST /api/v1/organizations` - Create organization
- `POST /api/v1/organizations/join` - Join organization via invite code
- `GET /api/v1/organizations` - Get user organizations
- `GET /api/v1/organizations/:orgId` - Get organization details
- `GET /api/v1/organizations/:orgId/members` - Get organization members
- `DELETE /api/v1/organizations/:orgId/members/:memberId` - Remove member
- `PUT /api/v1/organizations/:orgId` - Update organization

### Projects

- `POST /api/v1/organizations/:orgId/projects` - Create project
- `GET /api/v1/organizations/:orgId/projects` - Get organization projects
- `GET /api/v1/organizations/:orgId/projects/:projectId` - Get project details
- `PUT /api/v1/organizations/:orgId/projects/:projectId` - Update project
- `DELETE /api/v1/organizations/:orgId/projects/:projectId` - Delete project
- `PATCH /api/v1/organizations/:orgId/projects/:projectId/move` - Move project
- `PATCH /api/v1/organizations/:orgId/projects/move` - Bulk move projects

### Tasks

- `POST /api/v1/organizations/:orgId/projects/:projectId/tasks` - Create task
- `GET /api/v1/organizations/:orgId/projects/:projectId/tasks` - Get project tasks
- `GET /api/v1/organizations/:orgId/projects/:projectId/tasks/:taskId` - Get task details
- `PUT /api/v1/organizations/:orgId/projects/:projectId/tasks/:taskId` - Update task
- `DELETE /api/v1/organizations/:orgId/projects/:projectId/tasks/:taskId` - Delete task
- `PATCH /api/v1/organizations/:orgId/projects/:projectId/tasks/:taskId/move` - Move task
- `PATCH /api/v1/organizations/:orgId/projects/:projectId/tasks/move` - Bulk move tasks

### WebSocket

- `GET /api/v1/ws?org_id=:orgId` - WebSocket connection for real-time updates

## Key Changes from Gin Version

### Framework Migration
- **Gin → Fiber**: Switched to high-performance Fiber framework
- **Manual DB → GORM**: Replaced manual SQL with GORM ORM
- **Manual Migrations → Auto-migrations**: GORM handles schema automatically

### Database Layer
- **GORM Integration**: All models use GORM struct tags
- **Auto-migrations**: Database schema created automatically on startup
- **Relationships**: GORM handles foreign key relationships
- **Soft Deletes**: Built-in soft delete support

### Performance Improvements
- **Fiber Performance**: 3x faster than Gin in benchmarks
- **GORM Optimizations**: Efficient queries and relationships
- **Connection Pooling**: Optimized database connections
- **WebSocket**: Native Fiber WebSocket support

### Code Quality
- **Type Safety**: Strong typing with GORM models
- **Validation**: Built-in request validation
- **Error Handling**: Comprehensive error management
- **Clean Architecture**: Maintained separation of concerns

## Database Schema

The application uses GORM auto-migrations to create the following tables:

- `users` - User accounts with soft deletes
- `organizations` - Organizations with invite codes
- `org_members` - Organization membership with roles
- `projects` - Projects within organizations
- `project_assignees` - Project assignments
- `tasks` - Tasks within projects
- `task_assignees` - Task assignments

## WebSocket Events

The WebSocket connection supports the following message types:

- `task_moved` - Task status changed
- `project_moved` - Project status changed
- `task_created` - New task created
- `task_updated` - Task updated
- `task_deleted` - Task deleted
- `project_created` - New project created
- `project_updated` - Project updated
- `project_deleted` - Project deleted
- `user_joined` - User joined organization
- `user_left` - User left organization

## Security

- JWT-based authentication
- Password hashing with bcrypt
- CORS protection
- Input validation
- SQL injection prevention (GORM)
- Role-based access control

## Development

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
golangci-lint run
```

## Deployment

The application is designed to be deployed on Google Cloud Platform:

- **App Engine**: For initial deployment
- **GKE**: For containerized scaling
- **Cloud Functions**: For serverless operations
- **Cloud Endpoints**: For API management

## Performance Benefits

### Fiber Advantages
- **3x Faster**: Fiber is significantly faster than Gin
- **Lower Memory**: Reduced memory footprint
- **Better Concurrency**: Improved concurrent request handling
- **Native WebSocket**: Built-in WebSocket support

### GORM Advantages
- **Auto-migrations**: No manual SQL migrations needed
- **Type Safety**: Compile-time type checking
- **Relationships**: Automatic foreign key handling
- **Query Optimization**: Efficient database queries

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.