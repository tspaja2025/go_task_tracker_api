# Build Go Binary
FROM golang:1.24-alpine AS builder

# Install git and build tools
RUN apk add --no-cache git build-base

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -0 /task-tracker ./cmd/api/main.go

# Create Runtime Image
FROM alpine:3.19

WORKDIR /

# Copy the prebuild binary from the builder stage
COPY --from=builder /task-tracker /task-tracker

# Expose the port API will run on
EXPOSE 8080

# Command to run the executable
CMD ["/task-tracker"]
