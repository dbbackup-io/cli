package mongodb

import (
	"github.com/dbbackup-io/cli/cmd/shared"
	"github.com/spf13/cobra"
)

var MongoDBCmd = &cobra.Command{
	Use:   "mongodb",
	Short: "Restore MongoDB database",
	Long:  `Restore MongoDB database from cloud storage`,
}

func init() {
	// Create storage destination restore commands
	s3RestoreCmd := createS3RestoreCommand()
	gcsRestoreCmd := createGCSRestoreCommand()
	azureRestoreCmd := createAzureRestoreCommand()
	localRestoreCmd := createLocalRestoreCommand()

	MongoDBCmd.AddCommand(s3RestoreCmd)
	MongoDBCmd.AddCommand(gcsRestoreCmd)
	MongoDBCmd.AddCommand(azureRestoreCmd)
	MongoDBCmd.AddCommand(localRestoreCmd)
}

func createS3RestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "s3",
		Short: "Restore MongoDB database from S3",
		Long:  `Restore MongoDB database from AWS S3`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleMongoDBRestore(cmd, args, "s3")
		},
	}

	// Add restore-specific flags
	shared.AddMongoDBRestoreFlags(cmd)
	shared.AddS3RestoreFlags(cmd)

	return cmd
}

func createGCSRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gcs",
		Short: "Restore MongoDB database from Google Cloud Storage",
		Long:  `Restore MongoDB database from Google Cloud Storage`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleMongoDBRestore(cmd, args, "gcs")
		},
	}

	// Add restore-specific flags
	shared.AddMongoDBRestoreFlags(cmd)
	shared.AddGCSRestoreFlags(cmd)

	return cmd
}

func createAzureRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "azure",
		Short: "Restore MongoDB database from Azure Blob Storage",
		Long:  `Restore MongoDB database from Azure Blob Storage`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleMongoDBRestore(cmd, args, "azure")
		},
	}

	// Add restore-specific flags
	shared.AddMongoDBRestoreFlags(cmd)
	shared.AddAzureRestoreFlags(cmd)

	return cmd
}

func createLocalRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Restore MongoDB database from local storage",
		Long:  `Restore MongoDB database from local filesystem`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleLocalRestore(cmd, args, "MongoDB")
		},
	}

	// Add restore-specific flags
	shared.AddMongoDBRestoreFlags(cmd)
	shared.AddLocalRestoreFlags(cmd)

	return cmd
}
