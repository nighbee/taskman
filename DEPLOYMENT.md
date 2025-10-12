# TaskMan Docker Deployment Guide

This guide covers deploying the TaskMan application to Google Cloud App Engine using Docker containers.

## 📁 Project Structure

```
taskman/
├── backend/
│   ├── Dockerfile
│   ├── app.yaml
│   ├── .dockerignore
│   └── ... (Go source files)
├── frontend/
│   ├── Dockerfile
│   ├── app.yaml
│   ├── nginx.conf
│   ├── .dockerignore
│   └── ... (React source files)
├── docker-compose.yml
├── deploy.sh
└── DEPLOYMENT.md
```

## 🐳 Docker Configuration

### Backend Dockerfile
- **Base Image**: `golang:1.21-alpine` (build) → `alpine:latest` (runtime)
- **Port**: 8080
- **Health Check**: `/health` endpoint
- **Security**: Non-root user execution

### Frontend Dockerfile
- **Base Image**: `node:18-alpine` (build) → `nginx:alpine` (runtime)
- **Port**: 80
- **Static Files**: Served via Nginx
- **SPA Support**: Client-side routing handled

## 🚀 Local Development

### Using Docker Compose

1. **Start all services**:
   ```bash
   docker-compose up -d
   ```

2. **View logs**:
   ```bash
   docker-compose logs -f
   ```

3. **Stop services**:
   ```bash
   docker-compose down
   ```

4. **Rebuild services**:
   ```bash
   docker-compose up --build
   ```

### Services Available
- **Frontend**: http://localhost:3000
- **Backend**: http://localhost:8080
- **Database**: localhost:5432

## ☁️ Google Cloud App Engine Deployment

### Prerequisites

1. **Install Google Cloud CLI**:
   ```bash
   # Download from: https://cloud.google.com/sdk/docs/install
   ```

2. **Authenticate**:
   ```bash
   gcloud auth login
   gcloud auth application-default login
   ```

3. **Create a project**:
   ```bash
   gcloud projects create your-project-id
   gcloud config set project your-project-id
   ```

4. **Enable APIs**:
   ```bash
   gcloud services enable appengine.googleapis.com
   gcloud services enable sqladmin.googleapis.com
   gcloud services enable cloudbuild.googleapis.com
   ```

### Database Setup

1. **Create Cloud SQL instance**:
   ```bash
   gcloud sql instances create taskman-db \
     --database-version=POSTGRES_15 \
     --tier=db-f1-micro \
     --region=us-central1
   ```

2. **Create database**:
   ```bash
   gcloud sql databases create taskman --instance=taskman-db
   ```

3. **Set password**:
   ```bash
   gcloud sql users set-password postgres \
     --instance=taskman-db \
     --password=your-secure-password
   ```

### Configuration Updates

Before deploying, update these files:

#### backend/app.yaml
```yaml
env_variables:
  DB_HOST: "/cloudsql/PROJECT_ID:REGION:INSTANCE_NAME"
  DB_PASSWORD: "your-database-password"
  CORS_ALLOWED_ORIGINS: "https://frontend-dot-PROJECT_ID.appspot.com"
```

#### frontend/app.yaml
```yaml
env_variables:
  VITE_API_URL: "https://backend-dot-PROJECT_ID.appspot.com/api/v1"
```

### Deployment Commands

#### Option 1: Using the deployment script
```bash
# Update PROJECT_ID in deploy.sh first
./deploy.sh
```

#### Option 2: Manual deployment
```bash
# Deploy backend
cd backend
gcloud app deploy app.yaml

# Deploy frontend
cd ../frontend
gcloud app deploy app.yaml
```

### Post-Deployment

1. **Check service status**:
   ```bash
   gcloud app services list
   ```

2. **View logs**:
   ```bash
   gcloud app logs tail -s backend
   gcloud app logs tail -s frontend
   ```

3. **Access URLs**:
   - Backend: `https://backend-dot-PROJECT_ID.appspot.com`
   - Frontend: `https://frontend-dot-PROJECT_ID.appspot.com`

## 🔧 Environment Variables

### Backend Environment Variables
- `DB_HOST`: Database host (Cloud SQL connection string)
- `DB_PORT`: Database port (5432)
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_SSL_MODE`: SSL mode (require/disable)
- `JWT_SECRET`: JWT signing secret
- `JWT_EXPIRY`: JWT expiration time
- `SERVER_PORT`: Server port (8080)
- `SERVER_HOST`: Server host (0.0.0.0)
- `LOG_LEVEL`: Log level (info/debug)
- `CORS_ALLOWED_ORIGINS`: Allowed CORS origins
- `RUN_MIGRATIONS`: Run migrations on startup (true/false)

### Frontend Environment Variables
- `VITE_API_URL`: Backend API URL

## 🛠️ Troubleshooting

### Common Issues

1. **Build failures**:
   - Check Dockerfile syntax
   - Verify all dependencies are included
   - Check .dockerignore files

2. **Database connection issues**:
   - Verify Cloud SQL instance is running
   - Check connection string format
   - Ensure proper IAM permissions

3. **CORS errors**:
   - Update CORS_ALLOWED_ORIGINS with correct frontend URL
   - Check frontend VITE_API_URL configuration

4. **Health check failures**:
   - Verify health check endpoints are accessible
   - Check application startup time

### Useful Commands

```bash
# View service details
gcloud app services describe backend
gcloud app services describe frontend

# View versions
gcloud app versions list

# Rollback to previous version
gcloud app versions migrate VERSION_ID --service=backend

# Delete old versions
gcloud app versions delete VERSION_ID --service=backend
```

## 📊 Monitoring

1. **Cloud Console**: Monitor services in Google Cloud Console
2. **Logs**: View logs in Cloud Logging
3. **Metrics**: Monitor performance in Cloud Monitoring
4. **Health Checks**: Monitor service health status

## 🔒 Security Considerations

1. **Environment Variables**: Never commit secrets to version control
2. **Database**: Use Cloud SQL with proper access controls
3. **CORS**: Configure appropriate CORS policies
4. **HTTPS**: App Engine provides HTTPS by default
5. **IAM**: Use least privilege principle for service accounts

## 💰 Cost Optimization

1. **Instance Sizing**: Use appropriate instance sizes
2. **Auto-scaling**: Configure proper scaling policies
3. **Database**: Use appropriate Cloud SQL tiers
4. **Monitoring**: Set up billing alerts

## 📝 Next Steps

1. Set up CI/CD pipeline with Cloud Build
2. Configure custom domain
3. Set up monitoring and alerting
4. Implement backup strategies
5. Set up staging environment
