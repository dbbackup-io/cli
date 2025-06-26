package shared

import "github.com/spf13/cobra"

// S3Flags holds S3-specific flags
type S3Flags struct {
	Region    string
	Bucket    string
	Path      string
	AccessKey string
	SecretKey string
}

// GCSFlags holds Google Cloud Storage flags
type GCSFlags struct {
	ProjectID         string
	Bucket            string
	Path              string
	ServiceAccountKey string
}

// AzureFlags holds Azure Blob Storage flags
type AzureFlags struct {
	AccountName string
	AccountKey  string
	Container   string
	Path        string
}

// LocalFlags holds local storage flags
type LocalFlags struct {
	Directory string
}

// CommonFlags holds common backup flags
type CommonFlags struct {
	Compression string
}

// AddS3Flags adds S3 flags to a command
func AddS3Flags(cmd *cobra.Command, flags *S3Flags) {
	cmd.Flags().StringVar(&flags.Region, "region", "us-east-1", "AWS region")
	cmd.Flags().StringVar(&flags.Bucket, "bucket", "", "S3 bucket name (required)")
	cmd.Flags().StringVar(&flags.Path, "path", "", "S3 path prefix")
	cmd.Flags().StringVar(&flags.AccessKey, "aws-access-key", "", "AWS access key")
	cmd.Flags().StringVar(&flags.SecretKey, "aws-secret-key", "", "AWS secret key")

	cmd.MarkFlagRequired("bucket")
}

// AddGCSFlags adds Google Cloud Storage flags to a command
func AddGCSFlags(cmd *cobra.Command, flags *GCSFlags) {
	cmd.Flags().StringVar(&flags.ProjectID, "project-id", "", "GCS project ID (required)")
	cmd.Flags().StringVar(&flags.Bucket, "bucket", "", "GCS bucket name (required)")
	cmd.Flags().StringVar(&flags.Path, "path", "", "GCS path prefix")
	cmd.Flags().StringVar(&flags.ServiceAccountKey, "service-account-key", "", "Service account key JSON")

	cmd.MarkFlagRequired("project-id")
	cmd.MarkFlagRequired("bucket")
}

// AddAzureFlags adds Azure Blob Storage flags to a command
func AddAzureFlags(cmd *cobra.Command, flags *AzureFlags) {
	cmd.Flags().StringVar(&flags.AccountName, "account-name", "", "Azure storage account name (required)")
	cmd.Flags().StringVar(&flags.AccountKey, "account-key", "", "Azure storage account key (required)")
	cmd.Flags().StringVar(&flags.Container, "container", "", "Azure blob container name (required)")
	cmd.Flags().StringVar(&flags.Path, "path", "", "Azure blob path prefix")

	cmd.MarkFlagRequired("account-name")
	cmd.MarkFlagRequired("account-key")
	cmd.MarkFlagRequired("container")
}

// AddLocalFlags adds local storage flags to a command
func AddLocalFlags(cmd *cobra.Command, flags *LocalFlags) {
	cmd.Flags().StringVar(&flags.Directory, "directory", "./backups", "Local directory for backups")
}

// AddCommonFlags adds common backup flags to a command
func AddCommonFlags(cmd *cobra.Command, flags *CommonFlags) {
	cmd.Flags().StringVar(&flags.Compression, "compression", "gz", "Compression type (gz, none)")
}

// Restore-specific storage flags
func AddS3RestoreFlags(cmd *cobra.Command) {
	cmd.Flags().String("region", "us-east-1", "AWS region")
	cmd.Flags().String("bucket", "", "S3 bucket name (required)")
	cmd.Flags().String("aws-access-key", "", "AWS access key")
	cmd.Flags().String("aws-secret-key", "", "AWS secret key")

	cmd.MarkFlagRequired("bucket")
}

func AddGCSRestoreFlags(cmd *cobra.Command) {
	cmd.Flags().String("project-id", "", "GCS project ID (required)")
	cmd.Flags().String("bucket", "", "GCS bucket name (required)")
	cmd.Flags().String("service-account-key", "", "Service account key JSON")

	cmd.MarkFlagRequired("project-id")
	cmd.MarkFlagRequired("bucket")
}

func AddAzureRestoreFlags(cmd *cobra.Command) {
	cmd.Flags().String("account-name", "", "Azure storage account name (required)")
	cmd.Flags().String("account-key", "", "Azure storage account key (required)")
	cmd.Flags().String("container", "", "Azure blob container name (required)")

	cmd.MarkFlagRequired("account-name")
	cmd.MarkFlagRequired("account-key")
	cmd.MarkFlagRequired("container")
}

func AddLocalRestoreFlags(cmd *cobra.Command) {
	cmd.Flags().String("directory", "./backups", "Local directory containing backups")
}
