#!/bin/bash
# Quick test runner for Phase 1 Authentication tests
# Run this script to validate the authentication implementation

set -e  # Exit on first error

echo "=========================================="
echo "Phase 1 Authentication Test Suite"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to run tests and capture result
run_test() {
    local name=$1
    local command=$2

    echo -e "${YELLOW}Running: ${name}${NC}"

    if eval "$command"; then
        echo -e "${GREEN}✓ ${name} PASSED${NC}"
        echo ""
        return 0
    else
        echo -e "${RED}✗ ${name} FAILED${NC}"
        echo ""
        return 1
    fi
}

# Track failures
FAILED=0

# 1. Auth Core Tests
run_test "Password Hashing Tests" \
    "go test ./internal/auth/ -run 'Test(Hash|Verify)Password' -v" || FAILED=1

run_test "Session Token Tests" \
    "go test ./internal/auth/ -run 'Test(Generate|Hash).*Token' -v" || FAILED=1

run_test "Session Manager Tests" \
    "go test ./internal/auth/ -run 'Test(Create|Validate|Destroy).*Session' -v" || FAILED=1

run_test "Context Helper Tests" \
    "go test ./internal/auth/ -run 'TestContext' -v" || FAILED=1

run_test "Middleware Tests" \
    "go test ./internal/auth/ -run 'TestMiddleware' -v" || FAILED=1

# 2. Handler Tests
run_test "Login Handler Tests" \
    "go test ./internal/handler/auth/ -run 'TestLogin' -v" || FAILED=1

run_test "Register Handler Tests" \
    "go test ./internal/handler/auth/ -run 'TestRegister' -v" || FAILED=1

run_test "Logout Handler Tests" \
    "go test ./internal/handler/auth/ -run 'TestLogout' -v" || FAILED=1

# 3. Integration Tests
run_test "Auth Flow Integration Tests" \
    "go test ./tests/integration/ -run 'TestAuthFlow' -v -tags=integration" || FAILED=1

# 4. Race Detection
echo -e "${YELLOW}Running race detection tests...${NC}"
run_test "Race Detection" \
    "go test ./internal/auth/... -race -short" || FAILED=1

# 5. Coverage Report
echo -e "${YELLOW}Generating coverage report...${NC}"
go test ./internal/auth/... -coverprofile=coverage.out 2>/dev/null || true
if [ -f coverage.out ]; then
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo -e "${GREEN}Total Coverage: ${COVERAGE}${NC}"
    rm coverage.out
fi
echo ""

# Summary
echo "=========================================="
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ ALL TESTS PASSED${NC}"
    echo "=========================================="
    exit 0
else
    echo -e "${RED}✗ SOME TESTS FAILED${NC}"
    echo "=========================================="
    exit 1
fi
