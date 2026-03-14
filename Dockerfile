# =============================================================================
# Multi-stage Dockerfile for Skalkaho
# =============================================================================
# Stage 1: Build the Svelte estimate builder
# Stage 2: Build the Go binary (with Mage)
# Stage 3: Create minimal production image
# =============================================================================

# -----------------------------------------------------------------------------
# Frontend Build Stage
# -----------------------------------------------------------------------------
FROM node:22-alpine AS frontend

WORKDIR /app/ui

# Copy package files first for better layer caching
COPY ui/package.json ui/package-lock.json ./
RUN npm ci

# Copy frontend source and build
COPY ui/ ./
RUN npm run build

# -----------------------------------------------------------------------------
# Go Build Stage
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Install mage
RUN go install github.com/magefile/mage@latest

# Copy source code
COPY . .

# Copy frontend build output
COPY --from=frontend /app/static/estimate-builder ./static/estimate-builder/

# Build server binary with mage
RUN mage build

# -----------------------------------------------------------------------------
# Production Stage
# -----------------------------------------------------------------------------
FROM alpine:3.20

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata

# Create non-root user
RUN addgroup -g 1000 -S skalkaho && \
    adduser -u 1000 -S skalkaho -G skalkaho

WORKDIR /app

# Copy binaries and static files from builder
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/bin/seed /app/seed
COPY --from=builder /app/static /app/static

RUN chown -R skalkaho:skalkaho /app

# Switch to non-root user
USER skalkaho

# Environment defaults
ENV ADDR=:8080
ENV ENVIRONMENT=production

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/server"]
