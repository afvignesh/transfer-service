#!/bin/bash

# Start the database
echo "Starting database with docker-compose..."
docker-compose up -d

# Wait for database to be ready
echo "Waiting for database to be ready..."
until docker-compose exec -T db pg_isready -U user -d internal_transfer; do
    echo "Database is not ready yet. Waiting..."
    sleep 2
done

echo "Database is ready!"

echo "Checking Go dependencies..."
go mod tidy

# Run the Go application
echo "Starting Go application..."
go run cmd/main.go 