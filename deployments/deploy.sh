#!/bin/bash

# Build the Lambda function
./deployments/build-lambda.sh

# Deploy to AWS Lambda
echo "Deploying to AWS Lambda..."

aws lambda create-function \
    --function-name go-crud-lambda \
    --runtime go1.x \
    --role arn:aws:iam::YOUR_ACCOUNT_ID:role/lambda-role \
    --handler main \
    --zip-file fileb://build/lambda-function.zip \
    --environment Variables="{DB_HOST=your-db-host,DB_PORT=5432,DB_USER=postgres,DB_PASSWORD=your-password,DB_NAME=crud_db}" \
    --timeout 30 \
    --memory-size 128

echo "Deployment completed!"