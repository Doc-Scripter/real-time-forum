# Use the official Golang image from the Docker Hub
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files from the backend directory
COPY backend/go.mod backend/go.sum /app/backend/

# Download dependencies
# RUN go mod tidy

# Copy the backend directory content into the container, maintaining the structure
COPY backend/ /app/backend/

# Build the Go applicatdoion (the main.go file is inside /app/backend)
WORKDIR /app/backend
RUN go build -o forum

# Copy the frontend files into the container
COPY frontend/ /app/frontend/

# Set permissions for the frontend directory
RUN chmod -R 755 /app/frontend

# Expose the port that your application listens on
EXPOSE 9111

# Specify the command to run the application
ENTRYPOINT ["/app/backend/forum"]