#!/bin/bash

set -e

echo "Building for AWS Lambda..."

# Create build directory
mkdir -p build

# Clean previous builds
rm -f bootstrap
rm -f build/lambda-function.zip

# Build for Linux - output MUST be named 'bootstrap'
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/lambda/main.go

# Verify the binary was created and is executable
file bootstrap
ls -lh bootstrap
chmod +x bootstrap

# Create zip file with the bootstrap binary
zip -j build/lambda-function.zip bootstrap

# Clean up
rm bootstrap

echo "Lambda deployment package created: build/lambda-function.zip"
echo "Contents of zip file:"
unzip -l build/lambda-function.zip