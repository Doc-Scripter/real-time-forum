#!/bin/bash

# Build the Docker image
echo "Building Docker image..."
docker image prune -a 
docker build -t forum .

# Check if a docker instance is running
if [ $(docker ps -q -f name=forum)]; then
    echo "Stopping existing container..."
    docker stop forum
    docker rm forum
fi

# Run the Docker container from the image
echo "Running the Docker container..."
docker run forum
