package gcs

import (
	"context"
	"fmt"
	"io"
)

type Uploader struct {
	ProjectID         string
	Bucket            string
	ServiceAccountKey string
}

func (u *Uploader) Upload(ctx context.Context, key string, reader io.Reader) (int64, error) {
	// TODO: Implement Google Cloud Storage upload
	// This would use the Google Cloud Storage Go client library
	// For now, return an error indicating it's not implemented
	return 0, fmt.Errorf("google Cloud Storage upload not implemented yet")
}

// GetStorageType returns the storage type
func (u *Uploader) GetStorageType() string {
	return "gcs"
}
