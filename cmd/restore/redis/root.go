package redis

import (
	"github.com/dbbackup-io/cli/cmd/shared"
	"github.com/spf13/cobra"
)

var RedisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Restore Redis database",
	Long:  `Restore Redis database from cloud storage`,
}

func init() {
	// Create storage destination restore commands
	s3RestoreCmd := createS3RestoreCommand()
	gcsRestoreCmd := createGCSRestoreCommand()
	azureRestoreCmd := createAzureRestoreCommand()
	localRestoreCmd := createLocalRestoreCommand()

	RedisCmd.AddCommand(s3RestoreCmd)
	RedisCmd.AddCommand(gcsRestoreCmd)
	RedisCmd.AddCommand(azureRestoreCmd)
	RedisCmd.AddCommand(localRestoreCmd)
}

func createS3RestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "s3",
		Short: "Restore Redis database from S3",
		Long:  `Restore Redis database from AWS S3`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleRedisRestore(cmd, args, "s3")
		},
	}

	// Add restore-specific flags
	shared.AddRedisRestoreFlags(cmd)
	shared.AddS3RestoreFlags(cmd)

	return cmd
}

func createGCSRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gcs",
		Short: "Restore Redis database from Google Cloud Storage",
		Long:  `Restore Redis database from Google Cloud Storage`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleRedisRestore(cmd, args, "gcs")
		},
	}

	// Add restore-specific flags
	shared.AddRedisRestoreFlags(cmd)
	shared.AddGCSRestoreFlags(cmd)

	return cmd
}

func createAzureRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "azure",
		Short: "Restore Redis database from Azure Blob Storage",
		Long:  `Restore Redis database from Azure Blob Storage`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleRedisRestore(cmd, args, "azure")
		},
	}

	// Add restore-specific flags
	shared.AddRedisRestoreFlags(cmd)
	shared.AddAzureRestoreFlags(cmd)

	return cmd
}

func createLocalRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Restore Redis database from local storage",
		Long:  `Restore Redis database from local filesystem`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleLocalRestore(cmd, args, "Redis")
		},
	}

	// Add restore-specific flags
	shared.AddRedisRestoreFlags(cmd)
	shared.AddLocalRestoreFlags(cmd)

	return cmd
}
