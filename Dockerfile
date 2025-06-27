# Build stage
FROM golang:1.24-alpine AS builder

# Install necessary database clients
RUN apk add --no-cache \
    postgresql-client \
    mysql-client \
    mongodb-tools \
    redis

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary for target architecture
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o dbbackup main.go

# Final stage
FROM alpine:latest

# Install database clients in the final image
RUN apk add --no-cache \
    postgresql-client \
    mysql-client \
    mongodb-tools \
    redis \
    ca-certificates \
    tzdata

# Create non-root user
RUN addgroup -g 1001 -S dbbackup && \
    adduser -u 1001 -S dbbackup -G dbbackup

# Copy the binary from builder stage
COPY --from=builder /app/dbbackup /usr/local/bin/dbbackup

# Set permissions
RUN chmod +x /usr/local/bin/dbbackup

# Switch to non-root user
USER dbbackup

# Set working directory
WORKDIR /backups

# Default command
ENTRYPOINT ["dbbackup"]
CMD ["--help"]