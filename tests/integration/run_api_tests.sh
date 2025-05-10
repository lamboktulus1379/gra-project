#!/bin/bash

# Integration test script for GRA Framework migration
# This script starts the API server, runs tests against it, and reports results

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored text
print() {
  local color=$1
  local text=$2
  echo -e "${color}${text}${NC}"
}

# Start API server
start_server() {
  print $BLUE "Starting API server..."
  go run cmd/core-api/main.go &
  SERVER_PID=$!
  sleep 2 # Wait for server to start
  
  # Check if server is running
  if ! curl -s http://localhost:8080 > /dev/null; then
    print $RED "Failed to start server!"
    exit 1
  fi
  
  print $GREEN "Server started with PID: $SERVER_PID"
}

# Stop API server
stop_server() {
  print $BLUE "Stopping API server..."
  kill $SERVER_PID
  wait $SERVER_PID 2>/dev/null
  print $GREEN "Server stopped"
}

# Run tests
run_tests() {
  print $BLUE "Running API tests..."
  
  # Test 1: Hello World
  response=$(curl -s http://localhost:8080/)
  if [[ "$response" == *"success"* ]]; then
    print $GREEN "✓ Hello World test passed"
  else
    print $RED "✗ Hello World test failed"
    echo "Response: $response"
    tests_failed=true
  fi
  
  # Test 2: User Registration
  response=$(curl -s -X POST http://localhost:8080/users \
    -H "Content-Type: application/json" \
    -d '{"first_name":"John","last_name":"Doe","email":"john.doe@example.com","password":"securepassword123"}')
  if [[ "$response" == *"success"* ]] && [[ "$response" == *"User registered"* ]]; then
    print $GREEN "✓ User Registration test passed"
  else
    print $RED "✗ User Registration test failed"
    echo "Response: $response"
    tests_failed=true
  fi
  
  # Test 3: Invalid Registration
  response=$(curl -s -X POST http://localhost:8080/users \
    -H "Content-Type: application/json" \
    -d '{"first_name":"","email":"not-an-email","password":"123"}')
  if [[ "$response" == *"error"* ]] && [[ "$response" == *"required"* ]]; then
    print $GREEN "✓ Invalid Registration test passed"
  else
    print $RED "✗ Invalid Registration test failed"
    echo "Response: $response"
    tests_failed=true
  fi
  
  # Test 4: User Login (only if user registration succeeded)
  if [[ "$tests_failed" != true ]]; then
    response=$(curl -s -X POST http://localhost:8080/login \
      -H "Content-Type: application/json" \
      -d '{"email":"john.doe@example.com","password":"securepassword123"}')
    if [[ "$response" == *"token"* ]]; then
      print $GREEN "✓ User Login test passed"
      
      # Extract token for protected route test
      token=$(echo $response | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
      
      # Test 5: Protected Route
      if [[ -n "$token" ]]; then
        response=$(curl -s http://localhost:8080/api/protected \
          -H "Authorization: Bearer $token")
        if [[ "$response" == *"success"* ]]; then
          print $GREEN "✓ Protected Route test passed"
        else
          print $RED "✗ Protected Route test failed"
          echo "Response: $response"
          tests_failed=true
        fi
      fi
    else
      print $RED "✗ User Login test failed"
      echo "Response: $response"
      tests_failed=true
    fi
  fi
}

# Main execution
print $BLUE "===== API Integration Tests ====="

# Track if any tests failed
tests_failed=false

# Start the server, run tests, then stop server
start_server
run_tests
stop_server

# Report final result
if [[ "$tests_failed" == true ]]; then
  print $RED "✗ Some tests failed! Please review the output."
  exit 1
else
  print $GREEN "✓ All tests passed!"
  exit 0
fi
