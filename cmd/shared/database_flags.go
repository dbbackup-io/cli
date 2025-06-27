package shared

import "github.com/spf13/cobra"

// DatabaseFlags holds common database connection flags
type DatabaseFlags struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// AddPostgreSQLFlags adds PostgreSQL-specific flags to a command
func AddPostgreSQLFlags(cmd *cobra.Command, flags *DatabaseFlags) {
	cmd.Flags().StringVar(&flags.Host, "db-host", "localhost", "PostgreSQL host")
	cmd.Flags().IntVar(&flags.Port, "db-port", 5432, "PostgreSQL port")
	cmd.Flags().StringVar(&flags.Database, "db-name", "", "Database name (required)")
	cmd.Flags().StringVar(&flags.Username, "db-user", "", "Database username")
	cmd.Flags().StringVar(&flags.Password, "db-password", "", "Database password")

	_ = cmd.MarkFlagRequired("db-name")
}

// AddMySQLFlags adds MySQL-specific flags to a command
func AddMySQLFlags(cmd *cobra.Command, flags *DatabaseFlags) {
	cmd.Flags().StringVar(&flags.Host, "db-host", "localhost", "MySQL host")
	cmd.Flags().IntVar(&flags.Port, "db-port", 3306, "MySQL port")
	cmd.Flags().StringVar(&flags.Database, "db-name", "", "Database name (required)")
	cmd.Flags().StringVar(&flags.Username, "db-user", "", "Database username")
	cmd.Flags().StringVar(&flags.Password, "db-password", "", "Database password")

	_ = cmd.MarkFlagRequired("db-name")
}

// AddMongoDBFlags adds MongoDB-specific flags to a command
func AddMongoDBFlags(cmd *cobra.Command, flags *DatabaseFlags) {
	cmd.Flags().StringVar(&flags.Host, "db-host", "localhost", "MongoDB host")
	cmd.Flags().IntVar(&flags.Port, "db-port", 27017, "MongoDB port")
	cmd.Flags().StringVar(&flags.Database, "db-name", "", "Database name")
	cmd.Flags().StringVar(&flags.Username, "db-user", "", "Database username")
	cmd.Flags().StringVar(&flags.Password, "db-password", "", "Database password")
}

// AddRedisFlags adds Redis-specific flags to a command
func AddRedisFlags(cmd *cobra.Command, flags *DatabaseFlags) {
	cmd.Flags().StringVar(&flags.Host, "db-host", "localhost", "Redis host")
	cmd.Flags().IntVar(&flags.Port, "db-port", 6379, "Redis port")
	cmd.Flags().StringVar(&flags.Password, "db-password", "", "Redis password")
}

// Restore-specific flag functions
func AddPostgreSQLRestoreFlags(cmd *cobra.Command) {
	cmd.Flags().String("target-host", "localhost", "Target PostgreSQL host")
	cmd.Flags().Int("target-port", 5432, "Target PostgreSQL port")
	cmd.Flags().String("target-db", "", "Target database name (required)")
	cmd.Flags().String("target-user", "", "Target database username")
	cmd.Flags().String("target-password", "", "Target database password")
	cmd.Flags().String("backup-file", "", "Backup file path/key to restore (required)")

	_ = cmd.MarkFlagRequired("target-db")
	cmd.MarkFlagRequired("backup-file")
}

func AddMySQLRestoreFlags(cmd *cobra.Command) {
	cmd.Flags().String("target-host", "localhost", "Target MySQL host")
	cmd.Flags().Int("target-port", 3306, "Target MySQL port")
	cmd.Flags().String("target-db", "", "Target database name (required)")
	cmd.Flags().String("target-user", "", "Target database username")
	cmd.Flags().String("target-password", "", "Target database password")
	cmd.Flags().String("backup-file", "", "Backup file path/key to restore (required)")

	_ = cmd.MarkFlagRequired("target-db")
	cmd.MarkFlagRequired("backup-file")
}

func AddMongoDBRestoreFlags(cmd *cobra.Command) {
	cmd.Flags().String("target-host", "localhost", "Target MongoDB host")
	cmd.Flags().Int("target-port", 27017, "Target MongoDB port")
	cmd.Flags().String("target-db", "", "Target database name")
	cmd.Flags().String("target-user", "", "Target database username")
	cmd.Flags().String("target-password", "", "Target database password")
	cmd.Flags().String("backup-file", "", "Backup file path/key to restore (required)")

	cmd.MarkFlagRequired("backup-file")
}

func AddRedisRestoreFlags(cmd *cobra.Command) {
	cmd.Flags().String("target-host", "localhost", "Target Redis host")
	cmd.Flags().Int("target-port", 6379, "Target Redis port")
	cmd.Flags().String("target-password", "", "Target Redis password")
	cmd.Flags().String("backup-file", "", "Backup file path/key to restore (required)")

	cmd.MarkFlagRequired("backup-file")
}
