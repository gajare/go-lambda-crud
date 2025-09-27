#!/bin/bash

echo "Building for AWS Lambda..."

# Create build directory
mkdir -p build

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o build/main cmd/lambda/main.go

# Create zip file
zip -j build/lambda-function.zip build/main

echo "Lambda deployment package created: build/lambda-function.zip"