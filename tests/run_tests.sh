#!/bin/bash

echo "Running Transfer Service Tests..."
echo "=================================="

# Run account service tests
echo "Running Account Service Tests..."
go test ./tests/service -v -run "Test.*Account.*"

echo ""
echo "Running Transaction Service Validation Tests..."
go test ./tests/service -v -run "Test.*Validation.*"

echo ""
echo "All tests completed!" 