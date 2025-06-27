# Docker Usage

## Quick Start

```bash
# Pull the latest image
docker pull ghcr.io/hakan.kurtulus/dbbackup:latest

# Run with help
docker run --rm ghcr.io/hakan.kurtulus/dbbackup:latest --help

# Example: PostgreSQL backup to local storage
docker run --rm \
  -v $(pwd)/backups:/backups \
  --network host \
  ghcr.io/hakan.kurtulus/dbbackup:latest \
  dump postgres local \
  --db-host 127.0.0.1 \
  --db-port 5432 \
  --db-name mydb \
  --db-user postgres \
  --db-password mypassword \
  --directory /backups
```

## Available Images

### Multi-architecture Support
The Docker images support the following architectures:
- `linux/amd64` (x86_64)
- `linux/arm64` (aarch64)
- `linux/arm/v7` (armv7l)

### Tags
- `latest` - Latest stable release
- `v1.0.0` - Specific version tags
- `weekly-YYYYMMDD` - Weekly automated builds
- `main-<sha>` - Development builds from main branch

## Examples

### PostgreSQL Backup
```bash
docker run --rm \
  -v $(pwd)/backups:/backups \
  --network host \
  ghcr.io/hakan.kurtulus/dbbackup:latest \
  dump postgres local \
  --db-host localhost \
  --db-port 5432 \
  --db-name postgres \
  --db-user postgres \
  --db-password password \
  --directory /backups
```

### MySQL Backup with Debug Logging
```bash
docker run --rm \
  -v $(pwd)/backups:/backups \
  --network host \
  -e LOG_LEVEL=debug \
  ghcr.io/hakan.kurtulus/dbbackup:latest \
  dump mysql local \
  --db-host localhost \
  --db-port 3306 \
  --db-name mydb \
  --db-user root \
  --db-password password \
  --directory /backups
```

### MongoDB Backup to S3
```bash
docker run --rm \
  --network host \
  -e AWS_ACCESS_KEY_ID=your_access_key \
  -e AWS_SECRET_ACCESS_KEY=your_secret_key \
  ghcr.io/hakan.kurtulus/dbbackup:latest \
  dump mongodb s3 \
  --db-host localhost \
  --db-port 27017 \
  --db-name mydb \
  --db-user username \
  --db-password password \
  --bucket my-backup-bucket \
  --region us-east-1
```

### Redis Backup
```bash
docker run --rm \
  -v $(pwd)/backups:/backups \
  --network host \
  ghcr.io/hakan.kurtulus/dbbackup:latest \
  dump redis local \
  --db-host localhost \
  --db-port 6379 \
  --db-password password \
  --directory /backups
```

## Environment Variables

- `LOG_LEVEL=debug` - Enable debug logging
- `AWS_ACCESS_KEY_ID` - AWS access key for S3 backups
- `AWS_SECRET_ACCESS_KEY` - AWS secret key for S3 backups
- `GOOGLE_APPLICATION_CREDENTIALS` - Path to GCS service account JSON

## Volume Mounts

- `/backups` - Default directory for local backups
- `/config` - Configuration files (if needed)

## Network Considerations

When running in Docker, you may need to:

1. **Use host networking** (`--network host`) to access local databases
2. **Use custom networks** for containerized databases
3. **Use external hostnames** for remote databases

```bash
# For containerized databases on the same host
docker run --rm \
  -v $(pwd)/backups:/backups \
  --network host \
  ghcr.io/hakan.kurtulus/dbbackup:latest \
  dump postgres local \
  --db-host host.docker.internal \
  --db-port 5432 \
  # ... other options
```

## Building Locally

```bash
# Build for current architecture
docker build -t dbbackup:local .

# Build for multiple architectures
docker buildx build --platform linux/amd64,linux/arm64 -t dbbackup:multi .
```

## Security

The Docker image:
- Runs as non-root user (`dbbackup:dbbackup`)
- Includes only necessary database clients
- Uses Alpine Linux for minimal attack surface
- Supports read-only root filesystem