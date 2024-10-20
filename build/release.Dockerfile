# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the source code and build the application
COPY . .
RUN go build -o main cmd/server/main.go

# Stage 2: Create a lightweight image to run the Go binary
FROM alpine:latest
WORKDIR /app

# Ensure /app exists and has correct permissions
RUN mkdir -p /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main /app/main

# Set executable permission (if necessary)
RUN chmod +x /app/main

# Expose the port your app runs on (optional, adjust if needed)
EXPOSE 8080

# Set the entrypoint to the binary
ENTRYPOINT [ "/app/main" ]
