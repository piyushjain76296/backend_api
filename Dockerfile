FROM golang:alpine AS builder

# Install build dependencies (needed for sqlite3 cgo)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application (enable CGO for sqlite3)
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Install sqlite runtime dependencies
RUN apk add --no-cache sqlite-libs

# Copy binary from builder
COPY --from=builder /app/server .

# Copy frontend static files
COPY frontend/ ./frontend/

# Set environment variables
ENV PORT=8080
ENV DATABASE_PATH=./data/tickets.db

# Create data directory for sqlite
RUN mkdir -p ./data

# Expose port
EXPOSE 8080

# Command to run
CMD ["./server"]
