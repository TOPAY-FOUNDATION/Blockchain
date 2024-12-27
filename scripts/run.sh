#!/bin/bash

# Script to run the Go blockchain application

echo "Starting the Blockchain Application..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go and try again."
    exit 1
fi

# Run the application
go run cmd/main.go

if [ $? -eq 0 ]; then
    echo "Blockchain Application started successfully."
else
    echo "Failed to start the Blockchain Application."
    exit 1
fi
