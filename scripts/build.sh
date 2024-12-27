#!/bin/bash

# Script to build the Go blockchain application

echo "Building the Blockchain Application..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go and try again."
    exit 1
fi

# Clean previous builds
if [ -d "./bin" ]; then
    rm -rf ./bin
fi
mkdir -p ./bin

# Build the application
go build -o ./bin/blockchain main.go

if [ $? -eq 0 ]; then
    echo "Blockchain Application built successfully."
    echo "Executable is located at ./bin/blockchain"
else
    echo "Failed to build the Blockchain Application."
    exit 1
fi
