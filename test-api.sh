#!/bin/bash

echo "Testing API endpoints..."

BASE_URL="http://localhost:8080"

# Health check
echo "1. Health check:"
curl -s $BASE_URL/health | jq .

# Create user
echo -e "\n2. Create user:"
curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","age":30}' | jq .

# Get all users
echo -e "\n3. Get all users:"
curl -s $BASE_URL/users | jq .

# Get user by ID (replace 1 with actual user ID)
echo -e "\n4. Get user by ID:"
curl -s $BASE_URL/users/1 | jq .

# Update user
echo -e "\n5. Update user:"
curl -s -X PUT $BASE_URL/users/1 \
    -H "Content-Type: application/json" \
    -d '{"name":"John Updated","email":"john.updated@example.com","age":35}' | jq .

# Delete user
echo -e "\n6. Delete user:"
curl -s -X DELETE $BASE_URL/users/1 | jq .

echo -e "\nTesting completed!"