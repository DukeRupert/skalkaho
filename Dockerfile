# =============================================================================
# Multi-stage Dockerfile for Skalkaho
# =============================================================================
# Stage 1: Build the Go binary
# Stage 2: Create minimal production image
# =============================================================================

# -----------------------------------------------------------------------------
# Build Stage
# -----------------------------------------------------------------------------
FROM golang:1.23-alpine AS builder

# Install build dependencies for SQLite (CGO required for go-sqlite3)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with CGO enabled for SQLite
ARG VERSION=dev
ARG COMMIT=unknown
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT}" \
    -o /app/server \
    ./cmd/server

# -----------------------------------------------------------------------------
# Production Stage
# -----------------------------------------------------------------------------
FROM alpine:3.20

# Install runtime dependencies for SQLite
RUN apk add --no-cache \
    ca-certificates \
    sqlite-libs \
    tzdata

# Create non-root user
RUN addgroup -g 1000 -S skalkaho && \
    adduser -u 1000 -S skalkaho -G skalkaho

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/server

# Create data directory for SQLite database
RUN mkdir -p /app/data && chown -R skalkaho:skalkaho /app

# Switch to non-root user
USER skalkaho

# Environment defaults
ENV ADDR=:8080
ENV DATABASE_PATH=/app/data/quotes.db
ENV ENVIRONMENT=production

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/server"]
