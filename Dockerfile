# Stage 1: Build the binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Download dependencies first (caching layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o go-uns ./cmd/uns/main.go

# Stage 2: Create the minimal production image
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/go-uns .

# Command to run the executable
CMD ["./go-uns"]
