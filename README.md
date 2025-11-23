# TaskMan - Collaborative Task Management System

<div align="center">
  <img src="./frontend/logo.png" alt="TaskMan Logo" width="150" height="150">
  
  <p><strong>A powerful, real-time collaborative task management platform deployed on Google Cloud Platform</strong></p>
  
  [![Deployed on Google Cloud](https://img.shields.io/badge/Deployed%20on-Google%20Cloud-4285F4?logo=google-cloud&logoColor=white)](https://cloud.google.com)
  [![Backend](https://img.shields.io/badge/Backend-Go%201.21-00ADD8?logo=go&logoColor=white)](https://golang.org/)
  [![Frontend](https://img.shields.io/badge/Frontend-React%2018-61DAFB?logo=react&logoColor=white)](https://reactjs.org/)
  [![Database](https://img.shields.io/badge/Database-PostgreSQL-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
  [![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
</div>

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Architecture](#-architecture)
- [Technology Stack](#-technology-stack)
- [Google Cloud Deployment](#-google-cloud-deployment)
- [Getting Started](#-getting-started)
- [Project Structure](#-project-structure)
- [API Documentation](#-api-documentation)
- [Real-time Features](#-real-time-features)
- [Security](#-security)
- [Development](#-development)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🌟 Overview

**TaskMan** is a modern, full-stack collaborative task management system designed for teams who need real-time synchronization and powerful project management capabilities. Built with clean architecture principles and deployed on **Google Cloud Platform**, TaskMan offers enterprise-grade reliability, scalability, and performance.

### Why TaskMan?

- ✅ **Real-time Collaboration**: See updates instantly with WebSocket technology
- ✅ **Cloud-Native**: Fully deployed on Google Cloud Platform (App Engine)
- ✅ **Scalable Architecture**: Clean architecture with Go backend and React frontend
- ✅ **Enterprise Security**: JWT authentication, role-based access control, and encrypted data
- ✅ **Beautiful UI**: Modern, responsive design built with React and Tailwind CSS
- ✅ **Kanban Boards**: Intuitive drag-and-drop task management
- ✅ **Multi-Organization**: Support for multiple organizations and projects

---

## ✨ Features

### Core Features

- **🏢 Organization Management**
  - Create and manage multiple organizations
  - Invite team members with unique organization codes
  - Role-based access control (Admin/Member)
  - Organization-wide settings and permissions

- **📊 Project Management**
  - Create projects within organizations
  - Kanban-style board with customizable columns
  - Project assignment and team collaboration
  - Bulk operations support

- **✅ Task Management**
  - Create, update, and delete tasks
  - Nested subtasks support
  - Drag-and-drop functionality
  - Task assignments and priorities
  - Rich task descriptions and metadata

- **⚡ Real-time Collaboration**
  - Live task updates via WebSocket
  - User presence tracking (online/offline)
  - Real-time notifications for all changes
  - Instant synchronization across all clients

- **🔐 Authentication & Security**
  - Secure user registration and login
  - JWT-based authentication
  - Password encryption with bcrypt
  - CORS protection
  - SQL injection prevention

- **📱 Responsive Design**
  - Mobile-first approach
  - Beautiful UI with shadcn/ui components
  - Dark mode support
  - Accessible and user-friendly

---

## 🏗️ Architecture

TaskMan follows a clean, modern architecture pattern:

```
┌─────────────────────────────────────────────────────────────┐
│                     Google Cloud Platform                    │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                   App Engine (Frontend)                │  │
│  │  React + TypeScript + Vite + Tailwind CSS             │  │
│  │  - Responsive UI with shadcn/ui components            │  │
│  │  - WebSocket client for real-time updates            │  │
│  │  - State management with Zustand                     │  │
│  └────────────────────┬──────────────────────────────────┘  │
│                       │ HTTPS/WSS                            │
│  ┌────────────────────▼──────────────────────────────────┐  │
│  │                   App Engine (Backend)                 │  │
│  │  Go + Fiber + GORM                                    │  │
│  │  - RESTful API endpoints                              │  │
│  │  - WebSocket server for real-time features           │  │
│  │  - JWT authentication middleware                     │  │
│  │  - Clean architecture (handlers/services/models)    │  │
│  └────────────────────┬──────────────────────────────────┘  │
│                       │ Secure Connection                    │
│  ┌────────────────────▼──────────────────────────────────┐  │
│  │               Cloud SQL (PostgreSQL)                   │  │
│  │  - Managed PostgreSQL database                        │  │
│  │  - Automatic backups and replication                 │  │
│  │  - High availability and scalability                │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Clean Architecture Layers

**Backend (Go)**
- **Handlers Layer**: HTTP request/response handling and validation
- **Services Layer**: Business logic and database operations
- **Models Layer**: Data structures and DTOs
- **Database Layer**: Connection management and migrations
- **Middleware Layer**: Authentication, CORS, logging
- **Auth Layer**: JWT token management

**Frontend (React + TypeScript)**
- **Pages**: Route-based page components
- **Components**: Reusable UI components
- **Hooks**: Custom React hooks for API calls and state
- **Store**: Global state management with Zustand
- **Lib**: Utility functions and helpers

---

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.21
- **Web Framework**: Fiber v2 (high-performance HTTP framework)
- **WebSocket**: Fiber WebSocket support
- **Database ORM**: GORM (PostgreSQL)
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **Environment**: godotenv

### Frontend
- **Framework**: React 18
- **Language**: TypeScript 5
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **UI Components**: shadcn/ui (Radix UI)
- **State Management**: Zustand
- **HTTP Client**: Fetch API / TanStack Query
- **Drag & Drop**: @dnd-kit
- **Form Handling**: React Hook Form + Zod
- **Routing**: React Router DOM

### Database
- **Primary Database**: PostgreSQL 15
- **Hosting**: Google Cloud SQL
- **Features**: UUID primary keys, soft deletes, automatic timestamps, indexes

### Infrastructure
- **Cloud Provider**: Google Cloud Platform (GCP)
- **Compute**: App Engine (Flexible Environment)
- **Database**: Cloud SQL (PostgreSQL)
- **Containerization**: Docker
- **Orchestration**: Docker Compose (local development)
- **CI/CD**: Cloud Build

---

## ☁️ Google Cloud Deployment

**TaskMan is proudly deployed on Google Cloud Platform**, leveraging multiple GCP services for a robust, scalable, and reliable production environment.

### GCP Services Used

1. **App Engine (Flexible Environment)**
   - Hosts both frontend and backend services
   - Automatic scaling based on traffic
   - Built-in load balancing
   - Zero-downtime deployments
   - Custom domain support

2. **Cloud SQL (PostgreSQL)**
   - Managed PostgreSQL database
   - Automatic backups and point-in-time recovery
   - High availability with automatic failover
   - Secure private IP connectivity
   - Performance insights and monitoring

3. **Cloud Build**
   - Automated Docker image building
   - Containerized deployment pipeline
   - Fast and reliable builds

4. **Cloud Logging & Monitoring**
   - Centralized application logs
   - Performance metrics and alerting
   - Health check monitoring

### Deployment Architecture

```
Frontend (App Engine) → https://frontend-dot-[PROJECT_ID].appspot.com
Backend (App Engine)  → https://backend-dot-[PROJECT_ID].appspot.com
Database (Cloud SQL)  → Private network connection
```

### Why Google Cloud?

- ✅ **Scalability**: Automatic scaling based on demand
- ✅ **Reliability**: 99.95% uptime SLA
- ✅ **Security**: Enterprise-grade security and compliance
- ✅ **Performance**: Global CDN and edge locations
- ✅ **Cost-Effective**: Pay-only-for-what-you-use model
- ✅ **Managed Services**: No server management overhead

For detailed deployment instructions, see [DEPLOYMENT.md](./DEPLOYMENT.md)

---

## 🚀 Getting Started

### Prerequisites

- **Go 1.21+** (for backend development)
- **Node.js 18+** and npm (for frontend development)
- **PostgreSQL 15+** (for local database)
- **Docker & Docker Compose** (for containerized development)
- **Google Cloud CLI** (for deployment)

### Local Development Setup

#### Option 1: Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/nighbee/taskman.git
   cd taskman
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Backend Health: http://localhost:8080/health

4. **View logs**
   ```bash
   docker-compose logs -f
   ```

5. **Stop services**
   ```bash
   docker-compose down
   ```

#### Option 2: Manual Setup

**Backend Setup**

1. **Navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Run database migrations**
   ```bash
   make migrate-up
   # or
   go run main.go
   ```

5. **Start the backend server**
   ```bash
   make run
   # or
   go run main.go
   ```

**Frontend Setup**

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Configure environment variables**
   ```bash
   # Create .env file with:
   VITE_API_URL=http://localhost:8080/api/v1
   ```

4. **Start the development server**
   ```bash
   npm run dev
   ```

5. **Access the application**
   - Frontend: http://localhost:5173

---

## 📁 Project Structure

```
taskman/
├── backend/                    # Go backend service
│   ├── internal/              # Internal packages
│   │   ├── auth/             # JWT authentication
│   │   ├── config/           # Configuration management
│   │   ├── database/         # Database connection & migrations
│   │   ├── handlers/         # HTTP request handlers
│   │   ├── middleware/       # HTTP middleware
│   │   ├── models/           # Data models and DTOs
│   │   └── services/         # Business logic services
│   ├── migrations/           # Database migration files
│   ├── scripts/              # Utility scripts
│   ├── Dockerfile            # Backend Docker configuration
│   ├── app.yaml              # App Engine configuration
│   ├── go.mod                # Go module dependencies
│   └── main.go               # Application entry point
│
├── frontend/                  # React frontend application
│   ├── src/
│   │   ├── components/       # Reusable UI components
│   │   ├── hooks/            # Custom React hooks
│   │   ├── pages/            # Page components
│   │   ├── store/            # Zustand state management
│   │   ├── lib/              # Utility functions
│   │   ├── App.tsx           # Root component
│   │   └── main.tsx          # Application entry point
│   ├── public/               # Static assets
│   ├── Dockerfile            # Frontend Docker configuration
│   ├── nginx.conf            # Nginx configuration for production
│   ├── app.yaml              # App Engine configuration
│   ├── package.json          # npm dependencies
│   └── vite.config.ts        # Vite configuration
│
├── docker-compose.yml         # Docker Compose for local development
├── deploy.sh                  # Deployment script for GCP
├── DEPLOYMENT.md              # Detailed deployment guide
└── README.md                  # This file
```

---

## 📡 API Documentation

### Base URL
- **Local**: `http://localhost:8080/api/v1`
- **Production**: `https://backend-dot-[PROJECT_ID].appspot.com/api/v1`

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/register` | Register new user | No |
| POST | `/auth/login` | User login | No |
| GET | `/auth/me` | Get current user profile | Yes |
| POST | `/auth/refresh` | Refresh JWT token | Yes |

### Organization Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/orgs` | Create organization | Yes |
| GET | `/orgs` | List user's organizations | Yes |
| GET | `/orgs/:id` | Get organization details | Yes |
| PUT | `/orgs/:id` | Update organization | Yes (Admin) |
| DELETE | `/orgs/:id` | Delete organization | Yes (Admin) |
| POST | `/orgs/:id/join` | Join organization with code | Yes |
| GET | `/orgs/:id/members` | List organization members | Yes |

### Project Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/orgs/:orgId/projects` | Create project | Yes |
| GET | `/orgs/:orgId/projects` | List projects | Yes |
| GET | `/projects/:id` | Get project details | Yes |
| PUT | `/projects/:id` | Update project | Yes |
| DELETE | `/projects/:id` | Delete project | Yes |
| PATCH | `/projects/:id/move` | Move project status | Yes |
| GET | `/projects/:id/assignees` | List project assignees | Yes |

### Task Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/projects/:projectId/tasks` | Create task | Yes |
| GET | `/projects/:projectId/tasks` | List tasks | Yes |
| GET | `/tasks/:id` | Get task details | Yes |
| PUT | `/tasks/:id` | Update task | Yes |
| DELETE | `/tasks/:id` | Delete task | Yes |
| PATCH | `/tasks/:id/move` | Move task | Yes |
| GET | `/tasks/:id/assignees` | List task assignees | Yes |

### WebSocket Endpoint

| Endpoint | Description | Auth Required |
|----------|-------------|---------------|
| `/ws/:orgId` | WebSocket connection for real-time updates | Yes (via query param token) |

---

## ⚡ Real-time Features

TaskMan uses **WebSocket technology** for real-time collaboration:

### WebSocket Events

**Client → Server**
- `user_join`: User connects to organization
- `user_leave`: User disconnects from organization
- `ping`: Keep-alive heartbeat

**Server → Client**
- `task_created`: New task created
- `task_updated`: Task modified
- `task_deleted`: Task removed
- `task_moved`: Task position changed
- `project_created`: New project created
- `project_updated`: Project modified
- `project_deleted`: Project removed
- `project_moved`: Project status changed
- `user_joined`: New user joined organization
- `user_left`: User left organization
- `error`: Error notification

### Usage Example

```typescript
// Connect to WebSocket
const ws = new WebSocket(`ws://localhost:8080/ws/${orgId}?token=${jwtToken}`);

// Listen for events
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  switch(message.type) {
    case 'task_created':
      // Handle new task
      break;
    case 'task_updated':
      // Handle task update
      break;
    // ... handle other events
  }
};
```

---

## 🔒 Security

TaskMan implements multiple security layers:

### Authentication & Authorization
- **JWT Tokens**: Secure, stateless authentication
- **bcrypt**: Password hashing with salt
- **Role-based Access Control**: Admin and Member roles
- **Permission Checks**: Granular access control

### Data Protection
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: Parameterized queries with GORM
- **XSS Protection**: Input sanitization
- **HTTPS Enforcement**: TLS/SSL encryption
- **CORS Configuration**: Controlled cross-origin access

### Infrastructure Security
- **Cloud SQL**: Encrypted data at rest and in transit
- **Private Networking**: Secure backend-database communication
- **Environment Variables**: Secrets managed outside source code
- **IAM**: Least privilege principle for service accounts

---

## 💻 Development

### Backend Development

**Available Make Commands**
```bash
make run          # Run the backend server
make build        # Build the binary
make test         # Run tests
make migrate-up   # Run database migrations
make migrate-down # Rollback migrations
make docker-build # Build Docker image
make lint         # Run linter
```

**Database Migrations**
```bash
# Create new migration
go run scripts/create_migration.go migration_name

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down
```

### Frontend Development

**Available npm Scripts**
```bash
npm run dev       # Start development server
npm run build     # Build for production
npm run lint      # Run ESLint
npm run preview   # Preview production build
```

### Testing

**Backend Testing**
```bash
cd backend
go test ./... -v
```

**Frontend Testing**
```bash
cd frontend
npm run test
```

### Code Style

- **Backend**: Follow Go best practices and `gofmt` standards
- **Frontend**: ESLint configuration with React and TypeScript rules
- **Commits**: Use conventional commit messages

---

## 🤝 Contributing

We welcome contributions to TaskMan! Please follow these guidelines:

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Commit your changes** (`git commit -m 'Add amazing feature'`)
4. **Push to the branch** (`git push origin feature/amazing-feature`)
5. **Open a Pull Request**

### Development Guidelines

- Write clean, documented code
- Follow existing code style and architecture
- Add tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting PR

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Built with ❤️ using modern technologies
- Deployed on **Google Cloud Platform** for enterprise-grade reliability
- UI components from [shadcn/ui](https://ui.shadcn.com/)
- Icons from [Lucide](https://lucide.dev/)

---

## 📞 Support

For issues, questions, or contributions:
- **Issues**: [GitHub Issues](https://github.com/nighbee/taskman/issues)
- **Documentation**: See [DEPLOYMENT.md](./DEPLOYMENT.md) for deployment details

---

<div align="center">
  <p><strong>Built with Go, React, and deployed on Google Cloud Platform</strong></p>
  <p>⭐ Star this repository if you find it helpful!</p>
</div>
