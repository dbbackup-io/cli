package shared

import (
	"github.com/dbbackup-io/cli/pkg/backup"
	"github.com/spf13/cobra"
)

// DatabaseDumperFactory creates database dumpers from flags
type DatabaseDumperFactory func(flags DatabaseFlags) backup.DatabaseDumper

// CreateS3Command creates a generic S3 command for any database
func CreateS3Command(dbType string, dumperFactory DatabaseDumperFactory) *cobra.Command {
	var (
		dbFlags     DatabaseFlags
		s3Flags     S3Flags
		commonFlags CommonFlags
	)

	cmd := &cobra.Command{
		Use:   "s3",
		Short: "Export " + dbType + " database to S3",
		Long:  "Export " + dbType + " database backup to AWS S3",
		Run: func(cmd *cobra.Command, args []string) {
			dumper := dumperFactory(dbFlags)
			HandleS3Export(cmd, args, dumper, s3Flags, commonFlags)
		},
	}

	// Add appropriate flags based on database type
	switch dbType {
	case "PostgreSQL":
		AddPostgreSQLFlags(cmd, &dbFlags)
	case "MySQL":
		AddMySQLFlags(cmd, &dbFlags)
	case "MongoDB":
		AddMongoDBFlags(cmd, &dbFlags)
	case "Redis":
		AddRedisFlags(cmd, &dbFlags)
	}

	AddS3Flags(cmd, &s3Flags)
	AddCommonFlags(cmd, &commonFlags)

	return cmd
}

// CreateGCSCommand creates a generic GCS command for any database
func CreateGCSCommand(dbType string, dumperFactory DatabaseDumperFactory) *cobra.Command {
	var (
		dbFlags     DatabaseFlags
		gcsFlags    GCSFlags
		commonFlags CommonFlags
	)

	cmd := &cobra.Command{
		Use:   "gcs",
		Short: "Export " + dbType + " database to Google Cloud Storage",
		Long:  "Export " + dbType + " database backup to Google Cloud Storage",
		Run: func(cmd *cobra.Command, args []string) {
			dumper := dumperFactory(dbFlags)
			HandleGCSExport(cmd, args, dumper, gcsFlags, commonFlags)
		},
	}

	// Add appropriate flags based on database type
	switch dbType {
	case "PostgreSQL":
		AddPostgreSQLFlags(cmd, &dbFlags)
	case "MySQL":
		AddMySQLFlags(cmd, &dbFlags)
	case "MongoDB":
		AddMongoDBFlags(cmd, &dbFlags)
	case "Redis":
		AddRedisFlags(cmd, &dbFlags)
	}

	AddGCSFlags(cmd, &gcsFlags)
	AddCommonFlags(cmd, &commonFlags)

	return cmd
}

// CreateAzureCommand creates a generic Azure command for any database
func CreateAzureCommand(dbType string, dumperFactory DatabaseDumperFactory) *cobra.Command {
	var (
		dbFlags     DatabaseFlags
		azureFlags  AzureFlags
		commonFlags CommonFlags
	)

	cmd := &cobra.Command{
		Use:   "azure",
		Short: "Export " + dbType + " database to Azure Blob Storage",
		Long:  "Export " + dbType + " database backup to Azure Blob Storage",
		Run: func(cmd *cobra.Command, args []string) {
			dumper := dumperFactory(dbFlags)
			HandleAzureExport(cmd, args, dumper, azureFlags, commonFlags)
		},
	}

	// Add appropriate flags based on database type
	switch dbType {
	case "PostgreSQL":
		AddPostgreSQLFlags(cmd, &dbFlags)
	case "MySQL":
		AddMySQLFlags(cmd, &dbFlags)
	case "MongoDB":
		AddMongoDBFlags(cmd, &dbFlags)
	case "Redis":
		AddRedisFlags(cmd, &dbFlags)
	}

	AddAzureFlags(cmd, &azureFlags)
	AddCommonFlags(cmd, &commonFlags)

	return cmd
}

// CreateLocalCommand creates a generic local command for any database
func CreateLocalCommand(dbType string, dumperFactory DatabaseDumperFactory) *cobra.Command {
	var (
		dbFlags     DatabaseFlags
		localFlags  LocalFlags
		commonFlags CommonFlags
	)

	cmd := &cobra.Command{
		Use:   "local",
		Short: "Export " + dbType + " database to local storage",
		Long:  "Export " + dbType + " database backup to local filesystem",
		Run: func(cmd *cobra.Command, args []string) {
			dumper := dumperFactory(dbFlags)
			HandleLocalExport(cmd, args, dumper, localFlags, commonFlags)
		},
	}

	// Add appropriate flags based on database type
	switch dbType {
	case "PostgreSQL":
		AddPostgreSQLFlags(cmd, &dbFlags)
	case "MySQL":
		AddMySQLFlags(cmd, &dbFlags)
	case "MongoDB":
		AddMongoDBFlags(cmd, &dbFlags)
	case "Redis":
		AddRedisFlags(cmd, &dbFlags)
	}

	AddLocalFlags(cmd, &localFlags)
	AddCommonFlags(cmd, &commonFlags)

	return cmd
}
