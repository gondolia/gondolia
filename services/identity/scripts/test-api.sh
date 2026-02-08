#!/bin/bash

# Test API script for Identity Service
# Usage: ./scripts/test-api.sh [host]

HOST=${1:-http://localhost:8080}
TENANT="demo"

echo "Testing Identity Service API at $HOST"
echo "======================================="
echo ""

# Health check
echo "1. Health Check"
echo "   GET /health/live"
curl -s "$HOST/health/live" | jq .
echo ""

# Login as admin
echo "2. Login as admin"
echo "   POST /api/v1/auth/login"
RESPONSE=$(curl -s -X POST "$HOST/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT" \
  -d '{"email": "admin@demo.local", "password": "admin123"}')
echo "$RESPONSE" | jq .
ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r '.access_token')
echo ""

if [ "$ACCESS_TOKEN" == "null" ] || [ -z "$ACCESS_TOKEN" ]; then
  echo "Login failed. Make sure the database is seeded."
  exit 1
fi

# Get current user
echo "3. Get Current User (Me)"
echo "   GET /api/v1/auth/me"
curl -s "$HOST/api/v1/auth/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "X-Tenant-ID: $TENANT" | jq .
echo ""

# List users
echo "4. List Users"
echo "   GET /api/v1/users"
curl -s "$HOST/api/v1/users" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "X-Tenant-ID: $TENANT" | jq .
echo ""

# List companies
echo "5. List Companies"
echo "   GET /api/v1/companies"
curl -s "$HOST/api/v1/companies" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "X-Tenant-ID: $TENANT" | jq .
echo ""

# List roles
echo "6. List Roles"
echo "   GET /api/v1/roles"
curl -s "$HOST/api/v1/roles" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "X-Tenant-ID: $TENANT" | jq .
echo ""

# Logout
echo "7. Logout"
echo "   POST /api/v1/auth/logout"
REFRESH_TOKEN=$(echo "$RESPONSE" | jq -r '.refresh_token')
curl -s -X POST "$HOST/api/v1/auth/logout" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "X-Tenant-ID: $TENANT" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}" | jq .
echo ""

echo "======================================="
echo "API Tests completed!"
