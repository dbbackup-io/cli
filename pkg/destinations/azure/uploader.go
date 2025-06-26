package azure

import (
	"context"
	"fmt"
	"io"
)

type Uploader struct {
	AccountName string
	AccountKey  string
	Container   string
}

func (u *Uploader) Upload(ctx context.Context, key string, reader io.Reader) (int64, error) {
	// TODO: Implement Azure Blob Storage upload
	// This would use the Azure Blob Storage Go client library
	// For now, return an error indicating it's not implemented
	return 0, fmt.Errorf("Azure Blob Storage upload not implemented yet")
}

// GetStorageType returns the storage type
func (u *Uploader) GetStorageType() string {
	return "azure"
}