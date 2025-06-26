package dump

import (
	"github.com/dbbackup-io/cli/cmd/dump/mongodb"
	"github.com/dbbackup-io/cli/cmd/dump/mysql"
	"github.com/dbbackup-io/cli/cmd/dump/postgres"
	"github.com/dbbackup-io/cli/cmd/dump/redis"
	"github.com/spf13/cobra"
)

var DumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump database to storage",
	Long:  `Dump databases to cloud storage or local filesystem`,
}

func init() {
	// Add database subcommands to dump
	DumpCmd.AddCommand(postgres.PostgresCmd)
	DumpCmd.AddCommand(mysql.MySQLCmd)
	DumpCmd.AddCommand(mongodb.MongoDBCmd)
	DumpCmd.AddCommand(redis.RedisCmd)
}
