package redis

import (
	"github.com/dbbackup-io/cli/cmd/shared"
	"github.com/dbbackup-io/cli/pkg/backup"
	"github.com/dbbackup-io/cli/pkg/sources/redis"
	"github.com/spf13/cobra"
)

var RedisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Dump Redis database",
	Long:  `Dump Redis database to storage`,
}

// Redis dumper factory
func createRedisDumper(flags shared.DatabaseFlags) backup.DatabaseDumper {
	return &redis.Dumper{
		Host:     flags.Host,
		Port:     flags.Port,
		Password: flags.Password,
	}
}

func init() {
	// Create storage destination commands
	s3Cmd := shared.CreateS3Command("Redis", createRedisDumper)
	gcsCmd := shared.CreateGCSCommand("Redis", createRedisDumper)
	azureCmd := shared.CreateAzureCommand("Redis", createRedisDumper)
	localCmd := shared.CreateLocalCommand("Redis", createRedisDumper)

	RedisCmd.AddCommand(s3Cmd)
	RedisCmd.AddCommand(gcsCmd)
	RedisCmd.AddCommand(azureCmd)
	RedisCmd.AddCommand(localCmd)
}
