#!/bin/bash

# Script to deploy blockchain nodes using Docker Compose

echo "Deploying Blockchain Nodes..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker and try again."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "Docker Compose is not installed. Please install Docker Compose and try again."
    exit 1
fi

# Start the Docker Compose services
docker-compose up -d

if [ $? -eq 0 ]; then
    echo "Blockchain Nodes deployed successfully."
    echo "Run 'docker-compose ps' to check the status of the nodes."
else
    echo "Failed to deploy Blockchain Nodes."
    exit 1
fi
