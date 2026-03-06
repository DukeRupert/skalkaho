# =============================================================================
# Multi-stage Dockerfile for Skalkaho
# =============================================================================
# Stage 1: Build the Svelte frontend
# Stage 2: Build the Go binary
# Stage 3: Create minimal production image
# =============================================================================

# -----------------------------------------------------------------------------
# Frontend Build Stage
# -----------------------------------------------------------------------------
FROM node:22-alpine AS frontend

WORKDIR /app/frontend

# Copy package files first for better layer caching
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

# Copy frontend source and build
COPY frontend/ ./
RUN npm run build

# -----------------------------------------------------------------------------
# Go Build Stage
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

# Install build dependencies for SQLite (CGO required for go-sqlite3)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy frontend build output into static/js
COPY --from=frontend /app/static/js/quote-editor.iife.js ./static/js/quote-editor.iife.js

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

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    sqlite-libs \
    tzdata

# Create non-root user
RUN addgroup -g 1000 -S skalkaho && \
    adduser -u 1000 -S skalkaho -G skalkaho

WORKDIR /app

# Copy binary and static files from builder
COPY --from=builder /app/server /app/server
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
