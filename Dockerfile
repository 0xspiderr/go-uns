FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary from the new cmd directory
RUN CGO_ENABLED=0 GOOS=linux go build -o uns-app ./cmd/uns/main.go

# Stage 2: Create a minimal production image
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/uns-app .

# Command to run the executable
CMD ["./uns-app"]
