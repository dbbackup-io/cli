package restore

import (
	"github.com/dbbackup-io/cli/cmd/restore/mongodb"
	"github.com/dbbackup-io/cli/cmd/restore/mysql"
	"github.com/dbbackup-io/cli/cmd/restore/postgres"
	"github.com/dbbackup-io/cli/cmd/restore/redis"
	"github.com/spf13/cobra"
)

var RestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore database from cloud storage",
	Long:  `Restore databases from cloud storage backups`,
}

func init() {
	// Add database subcommands to restore
	RestoreCmd.AddCommand(postgres.PostgresCmd)
	RestoreCmd.AddCommand(mysql.MySQLCmd)
	RestoreCmd.AddCommand(mongodb.MongoDBCmd)
	RestoreCmd.AddCommand(redis.RedisCmd)
}
