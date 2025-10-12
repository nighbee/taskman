#!/bin/bash

# TaskMan Deployment Script for Google Cloud App Engine
# This script deploys both backend and frontend services

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ID="your-project-id"
REGION="us-central1"
DB_INSTANCE_NAME="taskman-db"

echo -e "${GREEN}🚀 Starting TaskMan deployment to Google Cloud App Engine${NC}"

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo -e "${RED}❌ gcloud CLI is not installed. Please install it first.${NC}"
    exit 1
fi

# Check if user is authenticated
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo -e "${YELLOW}⚠️  Please authenticate with gcloud first:${NC}"
    echo "gcloud auth login"
    exit 1
fi

# Set project
echo -e "${YELLOW}📋 Setting project to: ${PROJECT_ID}${NC}"
gcloud config set project $PROJECT_ID

# Enable required APIs
echo -e "${YELLOW}🔧 Enabling required APIs...${NC}"
gcloud services enable appengine.googleapis.com
gcloud services enable sqladmin.googleapis.com
gcloud services enable cloudbuild.googleapis.com

# Deploy backend
echo -e "${YELLOW}🔨 Deploying backend service...${NC}"
cd backend
gcloud app deploy app.yaml --quiet
cd ..

# Deploy frontend
echo -e "${YELLOW}🔨 Deploying frontend service...${NC}"
cd frontend
gcloud app deploy app.yaml --quiet
cd ..

echo -e "${GREEN}✅ Deployment completed successfully!${NC}"
echo -e "${GREEN}🌐 Backend URL: https://backend-dot-${PROJECT_ID}.appspot.com${NC}"
echo -e "${GREEN}🌐 Frontend URL: https://frontend-dot-${PROJECT_ID}.appspot.com${NC}"

echo -e "${YELLOW}📝 Next steps:${NC}"
echo "1. Update your Cloud SQL instance connection details in app.yaml"
echo "2. Update CORS_ALLOWED_ORIGINS with your frontend URL"
echo "3. Update VITE_API_URL in frontend/app.yaml with your backend URL"
echo "4. Run database migrations if needed"
