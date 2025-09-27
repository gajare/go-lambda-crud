# Go CRUD Microservice with AWS Lambda

A complete CRUD microservice built with Go, PostgreSQL, GORM, and deployed on AWS Lambda. This project demonstrates how to build a serverless application that can run both locally and on AWS Lambda.

## 📋 Table of Contents

- [Architecture Overview](#architecture-overview)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Local Development](#local-development)
- [AWS Lambda Deployment](#aws-lambda-deployment)
- [Testing with Curl](#testing-with-curl)
- [Terraform Deployment](#terraform-deployment)
- [Troubleshooting](#troubleshooting)

## 🏗️ Architecture Overview

### Dual Runtime Modes

This application supports two modes:

1. **Local HTTP Mode**: Runs as traditional web server using Gin framework
2. **AWS Lambda Mode**: Runs as serverless function when deployed to AWS

### Architecture Diagram

```
Local Development:
┌─────────────────┐    ┌─────────────────┐    ┌──────────────┐
│   HTTP Client   │ ── │   Go HTTP       │ ── │  PostgreSQL  │
│   (curl, etc.)  │    │   Server (Gin)  │    │  (Docker)    │
└─────────────────┘    └─────────────────┘    └──────────────┘

AWS Lambda Deployment:
┌─────────────────┐    ┌──────────────────┐    ┌──────────────┐
│   HTTP Client   │ ── │  API Gateway     │ ── │  AWS Lambda  │
│   (curl, etc.)  │    │  (AWS)           │    │  (Go Runtime)│
└─────────────────┘    └──────────────────┘    └───────┬──────┘
                                                       │
                                             ┌─────────▼────────┐
                                             │   Go Application │
                                             │  ┌─────────────┐ │
                                             │  │   Handler   │ │
                                             │  └─────────────┘ │
                                             │  ┌─────────────┐ │
                                             │  │  GORM ORM   │ │
                                             │  └─────────────┘ │
                                             └─────────┬────────┘
                                                       │
                                             ┌─────────▼────────┐
                                             │   PostgreSQL     │
                                             │     (RDS)        │
                                             └──────────────────┘
```

## 📁 Project Structure

```
go-lambda-crud/
├── cmd/
│   ├── lambda/          # Lambda entry point (bootstrap handler)
│   └── local/           # Local HTTP server entry point (Gin)
├── internal/
│   ├── handlers/        # Request handlers (Lambda + HTTP)
│   ├── models/          # Database models (GORM)
│   └── database/        # Database connection
├── pkg/
│   └── logger/          # Logging utilities (Logrus + Zap)
├── configs/             # Configuration management
├── deployments/         # Deployment scripts & Terraform
├── scripts/             # Utility scripts
├── docker-compose.yml   # Local development
├── Dockerfile          # Containerization
└── go.mod              # Go dependencies
```

## ⚙️ Prerequisites

### Required Software Installation

```bash
# 1. Install Go (1.21+)
# Download from: https://golang.org/dl/
go version  # Verify installation

# 2. Install Docker & Docker Compose
# Download from: https://www.docker.com/products/docker-desktop
docker --version && docker-compose --version  # Verify

# 3. Install AWS CLI
# Download from: https://aws.amazon.com/cli/
aws --version  # Verify

# 4. Install jq for JSON parsing (optional but helpful)
# macOS: brew install jq, Ubuntu: sudo apt install jq
```

### AWS Account Setup

```bash
# Configure AWS CLI with your credentials
aws configure
# Enter:
# - AWS Access Key ID
# - AWS Secret Access Key  
# - Default region (e.g., us-east-1)
# - Default output format (json)

# Verify AWS configuration
aws sts get-caller-identity
```

## 🚀 Local Development

### Step 1: Clone and Setup Project

```bash
# Clone the project
git clone <repository-url>
cd go-lambda-crud

# Install Go dependencies
go mod download

# Build the application
go build -o bin/local-main cmd/local/main.go
```

### Step 2: Start Local Services with Docker

```bash
# Start PostgreSQL and application
docker-compose up --build

# Or run in background
docker-compose up -d

# Check services are running
docker-compose ps

# View logs
docker-compose logs -f app
```

### Step 3: Verify Local Setup

```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {"status":"healthy"}
```

## 🌐 Testing with Curl (Local Development)

### Health Check
```bash
curl http://localhost:8080/health
```
**Response:**
```json
{"status":"healthy"}
```

### CRUD Operations - Local Testing

#### 1. Create a User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "age": 30
  }'
```

**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john.doe@example.com",
  "age": 30,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

#### 2. Get All Users
```bash
curl http://localhost:8080/users
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@example.com",
    "age": 30,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
]
```

#### 3. Get User by ID
```bash
curl http://localhost:8080/users/1
```

#### 4. Update User
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "john.smith@example.com",
    "age": 31
  }'
```

#### 5. Delete User
```bash
curl -X DELETE http://localhost:8080/users/1
```

**Response:**
```json
{"message": "User deleted successfully"}
```

### Automated Local Testing Script

```bash
# Make script executable
chmod +x scripts/test-local.sh

# Run comprehensive tests
./scripts/test-local.sh
```

## ☁️ AWS Lambda Deployment

### Step 1: Prepare for Lambda Deployment

```bash
# Build for AWS Lambda (Linux environment)
./deployments/build-lambda.sh

# Verify the build
unzip -l build/lambda-function.zip
# Should show: bootstrap (your Go binary)
```

### Step 2: Create IAM Role for Lambda

```bash
# Create execution role
./deployments/create-iam-role.sh

# Note the Role ARN output - you'll need it for deployment
```

### Step 3: Deploy to AWS Lambda

```bash
# Deploy using the complete script
./deployments/deploy-final.sh

# Or use one-command deployment
./deployments/deploy-one-click.sh
```

### Step 4: Verify Lambda Deployment

```bash
# Check function exists
aws lambda get-function --function-name go-crud-lambda

# Test function directly
aws lambda invoke \
    --function-name go-crud-lambda \
    --payload '{"httpMethod":"GET","path":"/health"}' \
    --cli-binary-format raw-in-base64-out \
    output.json && cat output.json
```

## 🌐 Testing with Curl (AWS Lambda + API Gateway)

### Step 1: Set Up API Gateway

```bash
# Set up HTTP API Gateway
./deployments/setup-api-gateway.sh

# Note the API URL output - this is your endpoint
```

### Step 2: Test Lambda via API Gateway

Replace `YOUR_API_URL` with the URL from the previous step.

#### Health Check
```bash
curl https://YOUR_API_ID.execute-api.YOUR_REGION.amazonaws.com/prod/health
```

#### Complete CRUD Testing

```bash
# Set your API URL
export API_URL="https://YOUR_API_ID.execute-api.YOUR_REGION.amazonaws.com/prod"

# 1. Health check
curl $API_URL/health

# 2. Create user
curl -X POST $API_URL/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Johnson","email":"alice@example.com","age":28}'

# 3. Get all users  
curl $API_URL/users

# 4. Get specific user (replace 1 with actual ID)
curl $API_URL/users/1

# 5. Update user
curl -X PUT $API_URL/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Smith","age":29}'

# 6. Delete user
curl -X DELETE $API_URL/users/1
```

### Automated Lambda Testing Script

```bash
# Create test script for Lambda
cat > test-lambda.sh << 'EOF'
#!/bin/bash
API_URL="https://YOUR_API_ID.execute-api.YOUR_REGION.amazonaws.com/prod"

echo "Testing Lambda API endpoints..."
echo "API URL: $API_URL"

echo -e "\n1. Health check:"
curl -s $API_URL/health | jq .

echo -e "\n2. Creating user..."
curl -s -X POST $API_URL/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","age":25}' | jq .

echo -e "\n3. Getting all users:"
curl -s $API_URL/users | jq .

echo -e "\n4. Test completed!"
EOF

chmod +x test-lambda.sh
./test-lambda.sh
```

## 🏗️ Terraform Deployment (Infrastructure as Code)

### Step 1: Initialize Terraform

```bash
cd deployments/terraform

# Initialize Terraform
terraform init

# Plan the deployment
terraform plan -var="db_host=localhost" -var="db_password=password"
```

### Step 2: Deploy with Terraform

```bash
# Apply the configuration
terraform apply -var="db_host=your-rds-endpoint" -var="db_password=your-password"

# Note the outputs (API URL, etc.)
```

### Step 3: Terraform Configuration Overview

The Terraform configuration creates:
- ✅ Lambda Function with proper IAM role
- ✅ API Gateway with REST API
- ✅ Resource permissions and integrations
- ✅ Environment variables configuration

### Step 4: Destroy Resources (When Done)

```bash
# Clean up all resources
terraform destroy
```

## 🔧 How Terraform Works Here

### Terraform File Structure
```
deployments/terraform/
├── main.tf           # Main configuration
├── variables.tf      # Input variables
├── outputs.tf        # Output values
└── terraform.tfstate # State file (generated)
```

### Key Terraform Components

#### 1. Lambda Function Resource
```hcl
resource "aws_lambda_function" "crud_lambda" {
  filename         = "../build/lambda-function.zip"
  function_name    = "go-crud-lambda"
  role            = aws_iam_role.lambda_role.arn
  handler         = "bootstrap"
  runtime         = "provided.al2"
  source_code_hash = filebase64sha256("../build/lambda-function.zip")
}
```

#### 2. API Gateway Integration
```hcl
resource "aws_api_gateway_rest_api" "crud_api" {
  name        = "go-crud-api"
  description = "Go CRUD API with Lambda"
}

resource "aws_api_gateway_integration" "integration" {
  rest_api_id             = aws_api_gateway_rest_api.crud_api.id
  resource_id             = aws_api_gateway_resource.users.id
  http_method             = aws_api_gateway_method.method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.crud_lambda.invoke_arn
}
```

## 📊 Performance Comparison

### Local Development
- **Startup Time**: Instant
- **Cost**: Free (local resources)
- **Database**: Local PostgreSQL in Docker
- **Best For**: Development and testing

### AWS Lambda
- **Startup Time**: Cold start ~100-300ms, warm start ~1-10ms
- **Cost**: Pay per request + compute time
- **Database**: Requires AWS RDS or similar
- **Best For**: Production, auto-scaling applications

## 🐛 Troubleshooting Guide

### Common Local Issues

#### Database Connection Issues
```bash
# Check if PostgreSQL is running
docker-compose ps

# View database logs
docker-compose logs postgres

# Reset database
docker-compose down -v
docker-compose up -d
```

#### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process or change port in docker-compose.yml
```

### Common Lambda Issues

#### Function Not Found
```bash
# Check if function exists
aws lambda list-functions

# Re-deploy if missing
./deployments/deploy-final.sh
```

#### Bootstrap Error
```bash
# Rebuild with correct bootstrap naming
rm -f bootstrap
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/lambda/main.go
zip -j build/lambda-function.zip bootstrap
```

#### Permission Issues
```bash
# Check IAM role
aws iam get-role --role-name lambda-go-crud-execution-role

# Update function role if needed
aws lambda update-function-configuration \
    --function-name go-crud-lambda \
    --role arn:aws:iam::YOUR_ACCOUNT:role/lambda-go-crud-execution-role
```

### Logging and Debugging

#### Local Logs
```bash
# View application logs
docker-compose logs -f app

# View database logs  
docker-compose logs -f postgres
```

#### Lambda Logs
```bash
# View Lambda logs
aws logs tail /aws/lambda/go-crud-lambda --follow

# Get recent logs
aws logs filter-log-events \
    --log-group-name /aws/lambda/go-crud-lambda \
    --limit 20
```

## 🚀 Quick Start Commands

### Local Development
```bash
# Start everything
docker-compose up --build

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/users

# Stop services
docker-compose down
```

### Lambda Deployment
```bash
# One-command deployment
./deployments/deploy-one-click.sh

# Test Lambda
aws lambda invoke --function-name go-crud-lambda \
    --payload '{"httpMethod":"GET","path":"/health"}' output.json

# Set up HTTP access
./deployments/setup-api-gateway.sh
```

### Complete Testing
```bash
# Local testing
./scripts/test-local.sh

# Lambda testing (after API Gateway setup)
./scripts/test-lambda.sh
```

## 📝 API Reference

### Endpoints Summary

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/users` | Get all users |
| POST | `/users` | Create new user |
| GET | `/users/{id}` | Get user by ID |
| PUT | `/users/{id}` | Update user |
| DELETE | `/users/{id}` | Delete user |

### Request/Response Examples

**Create User Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30
}
```

**Success Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

## 🎯 Key Features Implemented

- ✅ **Dual Runtime**: Local HTTP server + AWS Lambda
- ✅ **CRUD Operations**: Full user management
- ✅ **Database Integration**: PostgreSQL with GORM
- ✅ **Logging**: Dual logger (Logrus + Zap)
- ✅ **Containerization**: Docker for local development
- ✅ **Infrastructure as Code**: Terraform deployment
- ✅ **API Gateway**: HTTP access to Lambda
- ✅ **Auto-scaling**: Serverless architecture

## 🔄 Development Workflow

1. **Local Development**: Code and test locally with Docker
2. **Build for Lambda**: Use provided scripts to build deployment package
3. **Deploy**: Use AWS CLI or Terraform to deploy to Lambda
4. **Test**: Verify functionality via API Gateway
5. **Monitor**: Use CloudWatch logs for debugging

This setup provides a complete, production-ready serverless Go application that can scale efficiently while maintaining developer productivity through robust local development capabilities.