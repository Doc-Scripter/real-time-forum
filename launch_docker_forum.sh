#!/bin/bash

# Build the Docker image
echo "Building Docker image..."
docker image prune -a -f  # Remove unused images (use cautiously)
docker build -t forum .

# Check if a Docker instance is running or exists
if [ $(docker ps -a -q -f name=forum) ]; then
    echo "Stopping existing container..."
    docker stop forum
    docker rm forum
fi

# Run the Docker container from the image in detached mode
echo "Running the Docker container..."
docker run -d -p 9111:9111 --name forum forum