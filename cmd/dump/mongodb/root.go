package mongodb

import (
	"github.com/dbbackup-io/cli/cmd/shared"
	"github.com/dbbackup-io/cli/pkg/backup"
	"github.com/dbbackup-io/cli/pkg/sources/mongodb"
	"github.com/spf13/cobra"
)

var MongoDBCmd = &cobra.Command{
	Use:   "mongodb",
	Short: "Dump MongoDB database",
	Long:  `Dump MongoDB database to storage`,
}

// MongoDB dumper factory
func createMongoDBDumper(flags shared.DatabaseFlags) backup.DatabaseDumper {
	return &mongodb.Dumper{
		Host:     flags.Host,
		Port:     flags.Port,
		Database: flags.Database,
		Username: flags.Username,
		Password: flags.Password,
	}
}

func init() {
	// Create storage destination commands
	s3Cmd := shared.CreateS3Command("MongoDB", createMongoDBDumper)
	gcsCmd := shared.CreateGCSCommand("MongoDB", createMongoDBDumper)
	azureCmd := shared.CreateAzureCommand("MongoDB", createMongoDBDumper)
	localCmd := shared.CreateLocalCommand("MongoDB", createMongoDBDumper)

	MongoDBCmd.AddCommand(s3Cmd)
	MongoDBCmd.AddCommand(gcsCmd)
	MongoDBCmd.AddCommand(azureCmd)
	MongoDBCmd.AddCommand(localCmd)
}
