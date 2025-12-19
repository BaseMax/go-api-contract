#!/bin/bash

# Test script for go-api-contract CLI
# This script demonstrates all features of the CLI tool

set -e

echo "=================================="
echo "API Contract Validation Test Suite"
echo "=================================="
echo ""

# Build the tool
echo "1. Building the tool..."
go build -o apicontract ./cmd/apicontract
echo "✓ Build successful"
echo ""

# Start test server in the background
echo "2. Starting test server..."
go run examples/test-server.go &
SERVER_PID=$!
sleep 2
echo "✓ Test server started (PID: $SERVER_PID)"
echo ""

# Ensure server is killed on exit
trap "kill $SERVER_PID 2>/dev/null || true" EXIT

# Test 1: Basic GET request
echo "3. Test: Basic GET request"
./apicontract -spec examples/openapi-local.yaml \
  -endpoint http://localhost:8080/posts/1 \
  -method GET
echo ""

# Test 2: GET request with verbose mode
echo "4. Test: GET request with verbose mode"
./apicontract -spec examples/openapi-local.yaml \
  -endpoint http://localhost:8080/posts/1 \
  -method GET \
  -verbose
echo ""

# Test 3: JSON output
echo "5. Test: JSON output format"
./apicontract -spec examples/openapi-local.yaml \
  -endpoint http://localhost:8080/posts/1 \
  -method GET \
  -json
echo ""

# Test 4: POST request
echo "6. Test: POST request"
./apicontract -spec examples/openapi-local.yaml \
  -endpoint http://localhost:8080/posts \
  -method POST \
  -body '{"title":"Test Post","body":"Test content","userId":5}' \
  -verbose
echo ""

# Test 5: Environment variables
echo "7. Test: Environment variable substitution"
export API_BASE=http://localhost:8080
export POST_ID=2
./apicontract -spec examples/openapi-local.yaml \
  -endpoint '${API_BASE}/posts/${POST_ID}' \
  -method GET \
  -verbose
echo ""

# Test 6: Headers with environment variables
echo "8. Test: Headers with environment variables"
export API_TOKEN=my-secret-token
./apicontract -spec examples/openapi-local.yaml \
  -endpoint http://localhost:8080/posts/1 \
  -method GET \
  -headers '{"X-API-Key":"${API_TOKEN}","User-Agent":"apicontract/1.0"}' \
  -verbose
echo ""

# Test 7: 404 response validation
echo "9. Test: 404 response validation"
./apicontract -spec examples/openapi-local.yaml \
  -endpoint http://localhost:8080/posts/999 \
  -method GET \
  -verbose
echo ""

# Test 8: Get all posts
echo "10. Test: Get all posts (array response)"
./apicontract -spec examples/openapi-local.yaml \
  -endpoint http://localhost:8080/posts \
  -method GET \
  -json
echo ""

# Test 9: Version flag
echo "11. Test: Version flag"
./apicontract -version
echo ""

echo "=================================="
echo "✓ All tests passed!"
echo "=================================="

# Cleanup
kill $SERVER_PID 2>/dev/null || true
