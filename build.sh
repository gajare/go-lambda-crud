#!/bin/bash

set -e

echo "Building Lambda function for provided.al2 runtime..."

# Clean previous builds
rm -rf build
mkdir -p build

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/lambda/main.go

# Check file size and type
file bootstrap
ls -lh bootstrap

# Create zip with bootstrap file
zip -j build/lambda-function.zip bootstrap

# Clean up
rm bootstrap

echo "Build completed. File: build/lambda-function.zip"