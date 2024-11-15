# Use the official Golang image to build the app
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -v -o main .

# Expose the port the application listens on
EXPOSE 8080

# Command to run the application
CMD ["./main"]