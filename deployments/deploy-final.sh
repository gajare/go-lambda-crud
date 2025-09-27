#!/bin/bash

set -e

echo "=== Final Lambda Deployment ==="

# Configuration
FUNCTION_NAME="go-crud-lambda"
ROLE_ARN="arn:aws:iam::532678810964:role/lambda-go-crud-execution-role"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Step 1: Building with correct bootstrap naming...${NC}"

# Clean previous builds
rm -rf build
mkdir -p build
rm -f bootstrap

# Build with correct name
echo "Building Go binary as 'bootstrap'..."
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/lambda/main.go

# Verify build
if [ ! -f "bootstrap" ]; then
    echo -e "${RED}Error: bootstrap binary not created${NC}"
    exit 1
fi

echo "Binary details:"
file bootstrap
ls -lh bootstrap

# Make executable
chmod +x bootstrap

# Create zip
echo "Creating deployment package..."
zip -j build/lambda-function.zip bootstrap

# Verify zip contents
echo "ZIP contents:"
unzip -l build/lambda-function.zip

# Clean up
rm bootstrap

echo -e "${GREEN}✓ Build completed${NC}"

echo -e "${YELLOW}Step 2: Deploying to Lambda...${NC}"

# Update function code
aws lambda update-function-code \
    --function-name $FUNCTION_NAME \
    --zip-file fileb://build/lambda-function.zip

# Update configuration to ensure bootstrap handler
aws lambda update-function-configuration \
    --function-name $FUNCTION_NAME \
    --runtime provided.al2 \
    --handler bootstrap

echo -e "${GREEN}✓ Function updated${NC}"

echo -e "${YELLOW}Step 3: Waiting for update to complete...${NC}"
sleep 10

echo -e "${YELLOW}Step 4: Testing function...${NC}"

# Create better test event
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
        "requestId": "test-request",
        "accountId": "test",
        "stage": "test",
        "httpMethod": "GET",
        "path": "/health"
    }
}
EOF

# Invoke function with retry (sometimes it takes a moment to be ready)
for i in {1..3}; do
    echo "Attempt $i to invoke function..."
    
    aws lambda invoke \
        --function-name $FUNCTION_NAME \
        --payload file://test-event.json \
        --cli-binary-format raw-in-base64-out \
        --log-type Tail \
        output.json
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Function invoked successfully${NC}"
        break
    else
        echo "Attempt $i failed, retrying in 5 seconds..."
        sleep 5
    fi
done

# Show response
echo -e "\n${YELLOW}Response:${NC}"
cat output.json
echo ""

# Show logs with better formatting
echo -e "${YELLOW}Logs:${NC}"
aws lambda invoke \
    --function-name $FUNCTION_NAME \
    --payload file://test-event.json \
    --cli-binary-format raw-in-base64-out \
    --log-type Tail \
    output.json \
    --query 'LogResult' \
    --output text | base64 --decode

echo -e "\n${GREEN}=== Deployment Complete ===${NC}"
echo "Function: $FUNCTION_NAME"
echo "Runtime: provided.al2"
echo "Handler: bootstrap"

# Final verification
echo -e "\n${YELLOW}Final verification:${NC}"
aws lambda get-function --function-name $FUNCTION_NAME --query 'Configuration.{Runtime:Runtime, Handler:Handler, State:State}' --output table