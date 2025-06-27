package backup

import (
	"context"
	"io"
)

// DatabaseDumper interface for database backup sources
type DatabaseDumper interface {
	CreateBackupStream(ctx context.Context) (io.ReadCloser, error)
	GetFileExtension() string
	GetDatabaseType() string
}

// StorageUploader interface for storage destinations
type StorageUploader interface {
	Upload(ctx context.Context, key string, reader io.Reader) (int64, error)
	GetStorageType() string
}

// BackupConfig holds configuration for a backup operation
type BackupConfig struct {
	DatabaseType string
	DatabaseName string
	Compression  string
	PathPrefix   string
}

// BackupExecutor coordinates the backup process
type BackupExecutor struct {
	Dumper   DatabaseDumper
	Uploader StorageUploader
	Config   BackupConfig
}

func (be *BackupExecutor) Execute(ctx context.Context) error {
	// Generate filename
	filename := generateBackupFilename(be.Config, be.Dumper)

	fullPath := filename
	if be.Config.PathPrefix != "" {
		fullPath = be.Config.PathPrefix + "/" + filename
	}

	// Create backup stream
	reader, err := be.Dumper.CreateBackupStream(ctx)
	if err != nil {
		return err
	}

	var closeErr error
	defer func() {
		if err := reader.Close(); err != nil {
			closeErr = err
		}
	}()

	// Upload to storage
	size, err := be.Uploader.Upload(ctx, fullPath, reader)
	if err != nil {
		return err
	}

	// Check if the backup command itself failed
	if closeErr != nil {
		return closeErr
	}

	// Log success
	logBackupSuccess(fullPath, size, be.Dumper.GetDatabaseType(), be.Uploader.GetStorageType())
	return nil
}
