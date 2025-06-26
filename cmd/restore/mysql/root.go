package mysql

import (
	"github.com/dbbackup-io/cli/cmd/shared"
	"github.com/spf13/cobra"
)

var MySQLCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Restore MySQL database",
	Long:  `Restore MySQL database from cloud storage`,
}

func init() {
	// Create storage destination restore commands
	s3RestoreCmd := createS3RestoreCommand()
	gcsRestoreCmd := createGCSRestoreCommand()
	azureRestoreCmd := createAzureRestoreCommand()
	localRestoreCmd := createLocalRestoreCommand()

	MySQLCmd.AddCommand(s3RestoreCmd)
	MySQLCmd.AddCommand(gcsRestoreCmd)
	MySQLCmd.AddCommand(azureRestoreCmd)
	MySQLCmd.AddCommand(localRestoreCmd)
}

func createS3RestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "s3",
		Short: "Restore MySQL database from S3",
		Long:  `Restore MySQL database from AWS S3`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleMySQLRestore(cmd, args, "s3")
		},
	}

	// Add restore-specific flags
	shared.AddMySQLRestoreFlags(cmd)
	shared.AddS3RestoreFlags(cmd)

	return cmd
}

func createGCSRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gcs",
		Short: "Restore MySQL database from Google Cloud Storage",
		Long:  `Restore MySQL database from Google Cloud Storage`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleMySQLRestore(cmd, args, "gcs")
		},
	}

	// Add restore-specific flags
	shared.AddMySQLRestoreFlags(cmd)
	shared.AddGCSRestoreFlags(cmd)

	return cmd
}

func createAzureRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "azure",
		Short: "Restore MySQL database from Azure Blob Storage",
		Long:  `Restore MySQL database from Azure Blob Storage`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleMySQLRestore(cmd, args, "azure")
		},
	}

	// Add restore-specific flags
	shared.AddMySQLRestoreFlags(cmd)
	shared.AddAzureRestoreFlags(cmd)

	return cmd
}

func createLocalRestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Restore MySQL database from local storage",
		Long:  `Restore MySQL database from local filesystem`,
		Run: func(cmd *cobra.Command, args []string) {
			shared.HandleLocalRestore(cmd, args, "MySQL")
		},
	}

	// Add restore-specific flags
	shared.AddMySQLRestoreFlags(cmd)
	shared.AddLocalRestoreFlags(cmd)

	return cmd
}
