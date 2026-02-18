# Build stage
FROM golang:1.24-alpine AS builder
 
WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api
# Build migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate
# Build seeder tool
RUN CGO_ENABLED=0 GOOS=linux go build -o seeder ./cmd/seeder
# Build worker service
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/seeder .
COPY --from=builder /app/worker .
COPY --from=builder /app/.env .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
