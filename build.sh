#!/bin/bash

echo "Building Go Lambda CRUD application..."

# Build the application
go build -o bin/main cmd/lambda/main.go

echo "Build completed successfully!"
echo "To run locally: docker-compose up"