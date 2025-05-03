# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 1. Build stage
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Build the Go binary
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy dependency files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy project source code
COPY . .

# Build out binary (statically linked for Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 2. Runtime stage
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Minimal base image
FROM alpine:3.19

# Install security updates and CA certificates
RUN apk --no-cache upgrade && \
    apk --no-cache add ca-certificates

# Copy pre-built binary
WORKDIR /app
COPY --from=builder /app/server /app/server
# Set non-root user
RUN adduser -D -u 1001 appuser
USER appuser

# Expose port and run the server
EXPOSE 8080
CMD ["/app/server"]
