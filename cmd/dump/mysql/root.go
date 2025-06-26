package mysql

import (
	"github.com/dbbackup-io/cli/cmd/shared"
	"github.com/dbbackup-io/cli/pkg/backup"
	"github.com/dbbackup-io/cli/pkg/sources/mysql"
	"github.com/spf13/cobra"
)

var MySQLCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Dump MySQL database",
	Long:  `Dump MySQL database to storage`,
}

// MySQL dumper factory
func createMySQLDumper(flags shared.DatabaseFlags) backup.DatabaseDumper {
	return &mysql.Dumper{
		Host:     flags.Host,
		Port:     flags.Port,
		Database: flags.Database,
		Username: flags.Username,
		Password: flags.Password,
	}
}

func init() {
	// Create storage destination commands
	s3Cmd := shared.CreateS3Command("MySQL", createMySQLDumper)
	gcsCmd := shared.CreateGCSCommand("MySQL", createMySQLDumper)
	azureCmd := shared.CreateAzureCommand("MySQL", createMySQLDumper)
	localCmd := shared.CreateLocalCommand("MySQL", createMySQLDumper)

	MySQLCmd.AddCommand(s3Cmd)
	MySQLCmd.AddCommand(gcsCmd)
	MySQLCmd.AddCommand(azureCmd)
	MySQLCmd.AddCommand(localCmd)
}
