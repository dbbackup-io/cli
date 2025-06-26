package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Uploader struct {
	Directory string // Local directory to store backups
}

func (u *Uploader) Upload(ctx context.Context, key string, reader io.Reader) (int64, error) {
	// Ensure the directory exists
	if err := os.MkdirAll(u.Directory, 0755); err != nil {
		return 0, fmt.Errorf("failed to create directory %s: %w", u.Directory, err)
	}

	// Create the full file path
	filePath := filepath.Join(u.Directory, key)
	
	// Ensure parent directories exist
	parentDir := filepath.Dir(filePath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create parent directory %s: %w", parentDir, err)
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	// Copy data from reader to file
	size, err := io.Copy(file, reader)
	if err != nil {
		return 0, fmt.Errorf("failed to write backup data to %s: %w", filePath, err)
	}

	return size, nil
}

// GetStorageType returns the storage type
func (u *Uploader) GetStorageType() string {
	return "local"
}

// Download downloads a file from local storage to a writer
func (u *Uploader) Download(ctx context.Context, key string, writer io.Writer) error {
	filePath := filepath.Join(u.Directory, key)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", filePath)
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open backup file %s: %w", filePath, err)
	}
	defer file.Close()

	// Copy file content to writer
	_, err = io.Copy(writer, file)
	if err != nil {
		return fmt.Errorf("failed to read backup file %s: %w", filePath, err)
	}

	return nil
}

// GetFilePath returns the full file path for a backup key
func (u *Uploader) GetFilePath(key string) string {
	return filepath.Join(u.Directory, key)
}

// ListBackups lists all backup files in the directory
func (u *Uploader) ListBackups() ([]string, error) {
	var backups []string
	
	err := filepath.Walk(u.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			// Get relative path from directory
			relPath, err := filepath.Rel(u.Directory, path)
			if err != nil {
				return err
			}
			backups = append(backups, relPath)
		}
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to list backups in %s: %w", u.Directory, err)
	}
	
	return backups, nil
}