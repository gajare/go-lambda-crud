#!/bin/bash

echo "Starting local testing..."

# Start services
docker-compose up -d

# Wait for services to start
sleep 10

# Test endpoints
echo "Testing CREATE user..."
curl -X POST http://localhost:8080/users \
    -H "Content-Type: application/json" \
    -d '{"name":"John Doe","email":"john@example.com","age":30}'

echo -e "\nTesting GET all users..."
curl http://localhost:8080/users

echo -e "\nTesting GET user by ID..."
curl http://localhost:8080/users/1

echo -e "\nTesting UPDATE user..."
curl -X PUT http://localhost:8080/users/1 \
    -H "Content-Type: application/json" \
        -d '{"name":"John Updated","email":"john.updated@example.com","age":35}'

echo -e "\nTesting DELETE user..."
curl -X DELETE http://localhost:8080/users/1

echo -e "\nTesting completed!"