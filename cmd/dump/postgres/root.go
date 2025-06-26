package postgres

import (
	"github.com/dbbackup-io/cli/cmd/shared"
	"github.com/dbbackup-io/cli/pkg/backup"
	"github.com/dbbackup-io/cli/pkg/sources/postgres"
	"github.com/spf13/cobra"
)

var PostgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Dump PostgreSQL database",
	Long:  `Dump PostgreSQL database to storage`,
}

// PostgreSQL dumper factory
func createPostgresDumper(flags shared.DatabaseFlags) backup.DatabaseDumper {
	return &postgres.Dumper{
		Host:     flags.Host,
		Port:     flags.Port,
		Database: flags.Database,
		Username: flags.Username,
		Password: flags.Password,
	}
}

func init() {
	// Create storage destination commands
	s3Cmd := shared.CreateS3Command("PostgreSQL", createPostgresDumper)
	gcsCmd := shared.CreateGCSCommand("PostgreSQL", createPostgresDumper)
	azureCmd := shared.CreateAzureCommand("PostgreSQL", createPostgresDumper)
	localCmd := shared.CreateLocalCommand("PostgreSQL", createPostgresDumper)

	PostgresCmd.AddCommand(s3Cmd)
	PostgresCmd.AddCommand(gcsCmd)
	PostgresCmd.AddCommand(azureCmd)
	PostgresCmd.AddCommand(localCmd)
}
