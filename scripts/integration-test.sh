#!/bin/bash
# Gondolia Integration Test Script
# Tests all backend endpoints and validates responses
# No external dependencies required (jq-free version)

set -e  # Exit on first error

PROJECT_DIR="/home/juergen/projects/GoLandProjekte/gondolia"
BASE_URL="http://localhost:3001"
TENANT_ID="demo"
TEST_EMAIL="admin@demo.local"
TEST_PASSWORD="admin123"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
PASS_COUNT=0
FAIL_COUNT=0
TOTAL_COUNT=0

# Test results array
declare -a FAILED_TESTS

# Token storage
ACCESS_TOKEN=""

# Helper functions
pass() {
    echo -e "${GREEN}[PASS]${NC} $1"
    PASS_COUNT=$((PASS_COUNT + 1))
    TOTAL_COUNT=$((TOTAL_COUNT + 1))
}

fail() {
    echo -e "${RED}[FAIL]${NC} $1"
    FAILED_TESTS+=("$1")
    FAIL_COUNT=$((FAIL_COUNT + 1))
    TOTAL_COUNT=$((TOTAL_COUNT + 1))
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

info() {
    echo -e "${NC}[INFO]${NC} $1"
}

# Extract JSON field value (simple grep/sed approach)
get_json_field() {
    local json="$1"
    local field="$2"
    echo "$json" | grep -o "\"$field\"[[:space:]]*:[[:space:]]*\"[^\"]*\"" | sed "s/\"$field\"[[:space:]]*:[[:space:]]*\"\([^\"]*\)\"/\1/"
}

# Check if JSON contains field
has_json_field() {
    local json="$1"
    local field="$2"
    echo "$json" | grep -q "\"$field\""
}

# HTTP request helper with auth
api_get() {
    local path="$1"
    local headers=(-H "Content-Type: application/json" -H "X-Tenant-ID: $TENANT_ID")
    
    if [ -n "$ACCESS_TOKEN" ]; then
        headers+=(-H "Authorization: Bearer $ACCESS_TOKEN")
    fi
    
    curl -s -w "\n%{http_code}" "${headers[@]}" "${BASE_URL}${path}"
}

api_post() {
    local path="$1"
    local data="$2"
    local headers=(-H "Content-Type: application/json" -H "X-Tenant-ID: $TENANT_ID")
    
    if [ -n "$ACCESS_TOKEN" ]; then
        headers+=(-H "Authorization: Bearer $ACCESS_TOKEN")
    fi
    
    curl -s -w "\n%{http_code}" "${headers[@]}" -X POST -d "$data" "${BASE_URL}${path}"
}

# Extract HTTP status code from response
get_status() {
    echo "$1" | tail -n1
}

# Extract body from response
get_body() {
    echo "$1" | sed '$d'
}

echo "=================================================="
echo "  GONDOLIA INTEGRATION TEST SUITE"
echo "=================================================="
echo ""

# Prerequisite checks
info "Checking prerequisites..."
if ! command -v curl &> /dev/null; then
    echo -e "${RED}ERROR: curl not found${NC}"
    exit 1
fi

# 1. Docker Compose Health Check
info "Checking Docker Compose containers..."
RUNNING_CONTAINERS=$(cd "$PROJECT_DIR" && sg docker -c "docker compose ps" 2>/dev/null | grep -c "Up" || echo "0")
TOTAL_CONTAINERS=$(cd "$PROJECT_DIR" && sg docker -c "docker compose ps" 2>/dev/null | grep -v "NAME" | wc -l || echo "0")

if [ "$RUNNING_CONTAINERS" -eq "$TOTAL_CONTAINERS" ] && [ "$RUNNING_CONTAINERS" -gt 0 ]; then
    pass "All containers running ($RUNNING_CONTAINERS/$TOTAL_CONTAINERS)"
else
    fail "Some containers not running ($RUNNING_CONTAINERS/$TOTAL_CONTAINERS)"
fi

# 2. Health Endpoints
info "Testing health endpoints..."

RESPONSE=$(curl -s -w "\n%{http_code}" "${BASE_URL}/health/ready")
STATUS=$(get_status "$RESPONSE")
if [ "$STATUS" = "200" ]; then
    pass "Health check /health/ready"
else
    fail "Health check /health/ready (got $STATUS)"
fi

# 3. Login & Token
info "Testing authentication..."

LOGIN_DATA="{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}"
RESPONSE=$(api_post "/api/v1/auth/login" "$LOGIN_DATA")
STATUS=$(get_status "$RESPONSE")
BODY=$(get_body "$RESPONSE")

if [ "$STATUS" = "200" ]; then
    ACCESS_TOKEN=$(get_json_field "$BODY" "access_token")
    if [ -n "$ACCESS_TOKEN" ]; then
        pass "Login successful (token received)"
    else
        fail "Login response missing access_token"
        echo "Response: $BODY"
        exit 1
    fi
else
    fail "Login failed (status $STATUS)"
    echo "Response: $BODY"
    exit 1  # Can't continue without token
fi

# 4. Get User Info (/me)
info "Testing user endpoints..."

RESPONSE=$(api_get "/api/v1/auth/me")
STATUS=$(get_status "$RESPONSE")
BODY=$(get_body "$RESPONSE")

if [ "$STATUS" = "200" ]; then
    if has_json_field "$BODY" "user" && has_json_field "$BODY" "company"; then
        USER_EMAIL=$(get_json_field "$BODY" "email")
        pass "GET /api/v1/auth/me (user: ${USER_EMAIL:-authenticated})"
    else
        fail "GET /api/v1/auth/me - missing user or company fields"
        echo "Response: $BODY"
    fi
else
    fail "GET /api/v1/auth/me (status $STATUS)"
    echo "Response: $BODY"
fi

# 5. Catalog: Products List
info "Testing catalog endpoints..."

RESPONSE=$(api_get "/api/v1/products")
STATUS=$(get_status "$RESPONSE")
BODY=$(get_body "$RESPONSE")

if [ "$STATUS" = "200" ]; then
    if has_json_field "$BODY" "data"; then
        # Count items by counting commas in data array (rough estimate)
        PRODUCT_COUNT=$(echo "$BODY" | grep -o '"id"' | wc -l)
        if [ "$PRODUCT_COUNT" -gt 0 ]; then
            pass "GET /api/v1/products (~$PRODUCT_COUNT products)"
            
            # Store first product ID for detail test (extract first UUID-like pattern)
            FIRST_PRODUCT_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | sed 's/"id":"\([^"]*\)"/\1/')
        else
            fail "GET /api/v1/products - no products found"
        fi
    else
        fail "GET /api/v1/products - missing 'data' field"
        echo "Response: $BODY"
    fi
else
    fail "GET /api/v1/products (status $STATUS)"
    echo "Response: $BODY"
fi

# 6. Product Detail
if [ -n "$FIRST_PRODUCT_ID" ]; then
    RESPONSE=$(api_get "/api/v1/products/$FIRST_PRODUCT_ID")
    STATUS=$(get_status "$RESPONSE")
    BODY=$(get_body "$RESPONSE")
    
    if [ "$STATUS" = "200" ]; then
        if has_json_field "$BODY" "id" && has_json_field "$BODY" "sku"; then
            PRODUCT_SKU=$(get_json_field "$BODY" "sku")
            pass "GET /api/v1/products/:id (SKU: ${PRODUCT_SKU:-found})"
        else
            fail "GET /api/v1/products/:id - missing fields"
            echo "Response: $BODY"
        fi
    else
        fail "GET /api/v1/products/:id (status $STATUS)"
        echo "Response: $BODY"
    fi
fi

# 7. Product Prices
if [ -n "$FIRST_PRODUCT_ID" ]; then
    RESPONSE=$(api_get "/api/v1/products/$FIRST_PRODUCT_ID/prices")
    STATUS=$(get_status "$RESPONSE")
    BODY=$(get_body "$RESPONSE")
    
    if [ "$STATUS" = "200" ]; then
        if has_json_field "$BODY" "data"; then
            PRICE_COUNT=$(echo "$BODY" | grep -o '"price"' | wc -l)
            pass "GET /api/v1/products/:id/prices (~$PRICE_COUNT price scales)"
        else
            fail "GET /api/v1/products/:id/prices - missing 'data' field"
            echo "Response: $BODY"
        fi
    else
        fail "GET /api/v1/products/:id/prices (status $STATUS)"
        echo "Response: $BODY"
    fi
fi

# 8. Categories List
RESPONSE=$(api_get "/api/v1/categories")
STATUS=$(get_status "$RESPONSE")
BODY=$(get_body "$RESPONSE")

if [ "$STATUS" = "200" ]; then
    if has_json_field "$BODY" "data"; then
        CATEGORY_COUNT=$(echo "$BODY" | grep -o '"id"' | wc -l)
        if [ "$CATEGORY_COUNT" -gt 0 ]; then
            pass "GET /api/v1/categories (~$CATEGORY_COUNT categories)"
            
            # Store first category ID
            FIRST_CATEGORY_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | sed 's/"id":"\([^"]*\)"/\1/')
        else
            fail "GET /api/v1/categories - no categories found"
        fi
    else
        fail "GET /api/v1/categories - missing 'data' field"
        echo "Response: $BODY"
    fi
else
    fail "GET /api/v1/categories (status $STATUS)"
    echo "Response: $BODY"
fi

# 9. Category Detail
if [ -n "$FIRST_CATEGORY_ID" ]; then
    RESPONSE=$(api_get "/api/v1/categories/$FIRST_CATEGORY_ID")
    STATUS=$(get_status "$RESPONSE")
    BODY=$(get_body "$RESPONSE")
    
    if [ "$STATUS" = "200" ]; then
        if has_json_field "$BODY" "id" && has_json_field "$BODY" "name"; then
            pass "GET /api/v1/categories/:id (found)"
        else
            fail "GET /api/v1/categories/:id - missing fields"
            echo "Response: $BODY"
        fi
    else
        fail "GET /api/v1/categories/:id (status $STATUS)"
        echo "Response: $BODY"
    fi
fi

# 10. Database Data Validation
info "Validating database state..."

PRODUCT_DB_COUNT=$(cd "$PROJECT_DIR" && sg docker -c "docker exec gondolia-postgres-1 psql -U postgres -d catalog -t -c 'SELECT count(*) FROM products;'" 2>/dev/null | tr -d ' ' || echo "0")
if [ "$PRODUCT_DB_COUNT" -gt 0 ]; then
    pass "Database: $PRODUCT_DB_COUNT products in catalog DB"
else
    fail "Database: No products in catalog DB"
fi

CATEGORY_DB_COUNT=$(cd "$PROJECT_DIR" && sg docker -c "docker exec gondolia-postgres-1 psql -U postgres -d catalog -t -c 'SELECT count(*) FROM categories;'" 2>/dev/null | tr -d ' ' || echo "0")
if [ "$CATEGORY_DB_COUNT" -gt 0 ]; then
    pass "Database: $CATEGORY_DB_COUNT categories in catalog DB"
else
    fail "Database: No categories in catalog DB"
fi

TENANT_DB_COUNT=$(cd "$PROJECT_DIR" && sg docker -c "docker exec gondolia-postgres-1 psql -U postgres -d catalog -t -c \"SELECT count(*) FROM tenants WHERE code='demo';\"" 2>/dev/null | tr -d ' ' || echo "0")
if [ "$TENANT_DB_COUNT" -eq 1 ]; then
    pass "Database: Demo tenant exists in catalog DB"
else
    fail "Database: Demo tenant missing in catalog DB"
fi

# 11. Frontend API Path Validation
info "Validating frontend API client..."

CLIENT_FILE="$PROJECT_DIR/frontend/src/lib/api/client.ts"
if [ -f "$CLIENT_FILE" ]; then
    # Check for correct API paths (should NOT have /catalog prefix)
    if grep -q "/api/v1/catalog/" "$CLIENT_FILE"; then
        fail "Frontend API client uses wrong path /api/v1/catalog/* (should be /api/v1/products, /api/v1/categories)"
    else
        pass "Frontend API client uses correct paths"
    fi
    
    # Check for X-Tenant-ID header
    if grep -q "X-Tenant-ID" "$CLIENT_FILE"; then
        pass "Frontend API client sets X-Tenant-ID header"
    else
        fail "Frontend API client missing X-Tenant-ID header"
    fi
else
    warn "Frontend API client not found at $CLIENT_FILE"
fi

# 12. Docker Compose Config Validation
info "Validating docker-compose configuration..."

COMPOSE_FILE="$PROJECT_DIR/docker-compose.yml"
if grep -q "SECURE_COOKIES.*false" "$COMPOSE_FILE"; then
    pass "Docker Compose: SECURE_COOKIES=false (correct for HTTP)"
else
    fail "Docker Compose: SECURE_COOKIES should be 'false' for localhost HTTP"
fi

# Summary
echo ""
echo "=================================================="
echo "  TEST RESULTS SUMMARY"
echo "=================================================="
echo ""
echo "Total Tests:  $TOTAL_COUNT"
echo -e "${GREEN}Passed:       $PASS_COUNT${NC}"
echo -e "${RED}Failed:       $FAIL_COUNT${NC}"
echo ""

if [ ${#FAILED_TESTS[@]} -gt 0 ]; then
    echo "Failed Tests:"
    for test in "${FAILED_TESTS[@]}"; do
        echo -e "  ${RED}✗${NC} $test"
    done
    echo ""
fi

echo "=================================================="

# Exit with error if any test failed
if [ $FAIL_COUNT -gt 0 ]; then
    exit 1
else
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
fi
