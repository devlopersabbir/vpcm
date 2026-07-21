# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY cmd/ cmd/
COPY internal/ internal/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o vpsm-api ./cmd/vpsm-api

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/vpsm-api .

# Expose port 8080
EXPOSE 8080

# Command to run
CMD ["./vpsm-api"]
