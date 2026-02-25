# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files (from repo root)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY pkg/ ./pkg/
COPY provider/ ./provider/
COPY services/cart/ ./services/cart/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/service ./services/cart/cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy binaries from builder
COPY --from=builder /app/service .

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose ports
EXPOSE 8082 9092

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health/live || exit 1

# Run the application
ENTRYPOINT ["./service"]
