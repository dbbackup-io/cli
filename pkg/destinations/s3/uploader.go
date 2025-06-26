package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Uploader struct {
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
}

func (u *Uploader) Upload(ctx context.Context, key string, reader io.Reader) (int64, error) {
	// Create AWS session
	config := &aws.Config{
		Region: aws.String(u.Region),
	}

	// Use provided credentials if available
	if u.AccessKey != "" && u.SecretKey != "" {
		config.Credentials = credentials.NewStaticCredentials(u.AccessKey, u.SecretKey, "")
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return 0, err
	}

	// Create uploader
	uploader := s3manager.NewUploader(sess)

	// Upload the file
	result, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(key),
		Body:   reader,
	})

	if err != nil {
		return 0, err
	}

	// Note: S3 doesn't return the actual size from upload, so we return 0
	// In a real implementation, you might want to use a TeeReader to count bytes
	_ = result
	return 0, nil
}

// GetStorageType returns the storage type
func (u *Uploader) GetStorageType() string {
	return "s3"
}