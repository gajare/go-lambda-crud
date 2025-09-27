#!/bin/bash

set -e

echo "=== Deploying Go Lambda with Runtime Fix ==="

# Configuration
FUNCTION_NAME="go-crud-lambda"
ROLE_ARN="arn:aws:iam::532678810964:role/lambda-go-crud-execution-role"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Step 1: Building with bootstrap...${NC}"

# Build using the new script
chmod +x scripts/build-lambda.sh
./scripts/build-lambda.sh

if [ ! -f "build/lambda-function.zip" ]; then
    echo -e "${RED}Error: ZIP file not created${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Build completed${NC}"

echo -e "${YELLOW}Step 2: Checking existing function...${NC}"

# Check if function exists
EXISTS=$(aws lambda get-function --function-name $FUNCTION_NAME --query 'Configuration.FunctionName' --output text 2>/dev/null || echo "false")

if [ "$EXISTS" = "false" ]; then
    echo -e "${YELLOW}Creating new Lambda function...${NC}"
    
    aws lambda create-function \
        --function-name $FUNCTION_NAME \
        --runtime provided.al2 \
        --role "$ROLE_ARN" \
        --handler bootstrap \
        --zip-file fileb://build/lambda-function.zip \
        --environment Variables="{DB_HOST=localhost,DB_PORT=5432,DB_USER=postgres,DB_PASSWORD=password,DB_NAME=crud_db,LOG_LEVEL=info}" \
        --timeout 30 \
        --memory-size 256 \
        --description "Go CRUD Microservice Lambda Function"
    
    echo -e "${GREEN}✓ New function created${NC}"
else
    echo -e "${YELLOW}Updating existing function...${NC}"
    
    # Update function code
    aws lambda update-function-code \
        --function-name $FUNCTION_NAME \
        --zip-file fileb://build/lambda-function.zip
    
    # Update runtime configuration if needed
    aws lambda update-function-configuration \
        --function-name $FUNCTION_NAME \
        --runtime provided.al2 \
        --handler bootstrap \
        --environment Variables="{DB_HOST=localhost,DB_PORT=5432,DB_USER=postgres,DB_PASSWORD=password,DB_NAME=crud_db,LOG_LEVEL=info}"
    
    echo -e "${GREEN}✓ Function updated${NC}"
fi

echo -e "${YELLOW}Step 3: Testing function...${NC}"

# Wait for function to be active
echo "Waiting for function to be active..."
aws lambda wait function-updated --function-name $FUNCTION_NAME

# Create test event
cat > test-event.json << 'EOF'
{
    "httpMethod": "GET",
    "path": "/health",
    "headers": {
        "Content-Type": "application/json"
    },
    "queryStringParameters": null,
    "body": null,
    "pathParameters": null,
    "isBase64Encoded": false,
    "requestContext": {
        "requestId": "test-request-id"
    }
}
EOF

# Invoke function
echo "Testing function invocation..."
aws lambda invoke \
    --function-name $FUNCTION_NAME \
    --payload file://test-event.json \
    --cli-binary-format raw-in-base64-out \
    --log-type Tail \
    output.json

echo -e "${GREEN}✓ Function invoked${NC}"

# Show response
echo -e "\n${YELLOW}Response:${NC}"
cat output.json
echo ""

# Show logs
echo -e "${YELLOW}Logs:${NC}"
aws lambda invoke \
    --function-name $FUNCTION_NAME \
    --payload file://test-event.json \
    --cli-binary-format raw-in-base64-out \
    --log-type Tail \
    output.json \
    --query 'LogResult' \
    --output text | base64 --decode

echo -e "\n${GREEN}=== Deployment Successful ===${NC}"
echo "Function: $FUNCTION_NAME"
echo "Runtime: provided.al2"
echo "Handler: bootstrap"