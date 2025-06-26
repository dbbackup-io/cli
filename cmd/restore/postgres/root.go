package postgres

import (
	"github.com/dbbackup-io/cli/cmd/shared"
	"github.com/spf13/cobra"
)

var PostgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Restore PostgreSQL database",
	Long:  `Restore PostgreSQL database from cloud storage`,
}

func init() {
	// Create storage destination restore commands
	s3RestoreCmd := createS3RestoreCommand()
	gcsRestoreCmd := createGCSRestoreCommand()
	azureRestoreCmd := createAzureRestoreCommand()
	localRestoreCmd := createLocalRestoreCommand()

	PostgresCmd.AddCommand(s3RestoreCmd)
	PostgresCmd.AddCommand(gcsRestoreCmd)
	PostgresCmd.AddCommand(azureRestoreCmd)
	PostgresCmd.AddCommand(localRestoreCmd)
}

func createS3RestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "s3",
		Short: "Restore PostgreSQL database from S3",
		Long:  `Restore PostgreSQL database from AWS S3`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandlePostgresRestore(cmd, args, "s3")
		},
	}

	// Add restore-specific flags
	shared.AddPostgreSQLRestoreFlags(cmd)
	shared.AddS3RestoreFlags(cmd)

	return cmd
}

func createGCSRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gcs",
		Short: "Restore PostgreSQL database from Google Cloud Storage",
		Long:  `Restore PostgreSQL database from Google Cloud Storage`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandlePostgresRestore(cmd, args, "gcs")
		},
	}

	// Add restore-specific flags
	shared.AddPostgreSQLRestoreFlags(cmd)
	shared.AddGCSRestoreFlags(cmd)

	return cmd
}

func createAzureRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "azure",
		Short: "Restore PostgreSQL database from Azure Blob Storage",
		Long:  `Restore PostgreSQL database from Azure Blob Storage`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandlePostgresRestore(cmd, args, "azure")
		},
	}

	// Add restore-specific flags
	shared.AddPostgreSQLRestoreFlags(cmd)
	shared.AddAzureRestoreFlags(cmd)

	return cmd
}

func createLocalRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Restore PostgreSQL database from local storage",
		Long:  `Restore PostgreSQL database from local filesystem`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleLocalRestore(cmd, args, "PostgreSQL")
		},
	}

	// Add restore-specific flags
	shared.AddPostgreSQLRestoreFlags(cmd)
	shared.AddLocalRestoreFlags(cmd)

	return cmd
}
