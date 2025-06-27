package cmd

import (
	"fmt"
	"os"

	"github.com/dbbackup-io/cli/cmd/database_source"
	"github.com/dbbackup-io/cli/cmd/dump"
	"github.com/dbbackup-io/cli/cmd/job"
	"github.com/dbbackup-io/cli/cmd/login"
	"github.com/dbbackup-io/cli/cmd/logout"
	"github.com/dbbackup-io/cli/cmd/restore"
	"github.com/dbbackup-io/cli/cmd/server"
	"github.com/dbbackup-io/cli/cmd/status"
	"github.com/dbbackup-io/cli/cmd/storage_destination"
	"github.com/dbbackup-io/cli/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	logLevel      string
	logsTableView bool
	version       = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "dbbackup",
	Short:   "Database backup and restore CLI tool",
	Version: version,
	Long: `A CLI tool to backup and restore databases to/from cloud storage.
Supports PostgreSQL, MySQL, MongoDB, Redis with S3, GCS, and Azure storage.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set log level before any command runs
		if err := logger.SetLevel(logLevel); err != nil {
			fmt.Printf("Invalid log level: %v\n", err)
			os.Exit(1)
		}
	},
}

// Standalone logs command
var logsCmd = &cobra.Command{
	Use:   "logs [job-id]",
	Short: "View backup runs for a job",
	Long:  "List backup runs for a job. If no job ID is provided, you'll be prompted to select one interactively.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create args for job logs command
		jobArgs := []string{"logs"}

		// Add job ID if provided
		if len(args) > 0 {
			jobArgs = append(jobArgs, args[0])
		}

		// Add table flag if set
		tableFlag, _ := cmd.Flags().GetBool("table")
		if tableFlag {
			jobArgs = append(jobArgs, "--table")
		}

		// Pass control to job logs command
		job.JobCmd.SetArgs(jobArgs)
		job.JobCmd.Execute()
	},
}

// SetVersion sets the version for the CLI
func SetVersion(v string) {
	version = v
	rootCmd.Version = v
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set the logging level (error, warn, info, debug)")

	// Configure logs command
	logsCmd.Flags().Bool("table", false, "Display backup runs in table format instead of interactive selection")

	// Add all commands to root
	rootCmd.AddCommand(dump.DumpCmd)
	rootCmd.AddCommand(restore.RestoreCmd)
	rootCmd.AddCommand(login.LoginCmd)
	rootCmd.AddCommand(logout.LogoutCmd)
	rootCmd.AddCommand(status.StatusCmd)
	rootCmd.AddCommand(job.JobCmd)
	rootCmd.AddCommand(server.ServerCmd)
	rootCmd.AddCommand(database_source.DatabaseSourceCmd)
	rootCmd.AddCommand(storage_destination.StorageDestinationCmd)
	rootCmd.AddCommand(logsCmd)
}
