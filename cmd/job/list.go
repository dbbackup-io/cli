package job

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/dbbackup-io/cli/pkg/api"
	"github.com/dbbackup-io/cli/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	tableView bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List backup jobs",
	Long:  "List all backup jobs for the authenticated team. Use interactive selection to view job details.",
	Run:   runList,
}

func init() {
	listCmd.Flags().BoolVar(&tableView, "table", false, "Display jobs in table format instead of interactive selection")
	JobCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	logger.Debug("Starting backup jobs list")

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		logger.Errorf("Failed to create API client: %v", err)
		return
	}

	// Get backup jobs
	jobs, err := client.GetBackupJobs()
	if err != nil {
		logger.Errorf("Failed to retrieve backup jobs: %v", err)
		return
	}

	if len(jobs) == 0 {
		logger.Info("No backup jobs found")
		return
	}

	// Show table view if requested
	if tableView {
		displayJobsTable(jobs)
		return
	}

	// Interactive selection
	if err := selectAndDisplayJob(client, jobs); err != nil {
		logger.Debugf("Interactive selection failed: %v", err)
		logger.Info("Falling back to table view...")
		fmt.Println()
		displayJobsTable(jobs)
	}
}

// displayJobsTable shows jobs in a table format
func displayJobsTable(jobs []api.BackupJob) {
	fmt.Printf("%-5s %-25s %-10s %-15s %-15s %-20s\n", "ID", "NAME", "STATUS", "SCHEDULE", "MODE", "CREATED")
	fmt.Println(strings.Repeat("-", 90))

	for _, job := range jobs {
		// Truncate long names for display
		id := fmt.Sprintf("%d", job.ID)
		name := truncateString(job.Name, 25)
		status := truncateString(job.Status, 10)
		schedule := truncateString(job.Schedule, 15)
		mode := truncateString(job.ExecutionMode, 15)
		created := job.CreatedAt.Format("2006-01-02 15:04")

		fmt.Printf("%-5s %-25s %-10s %-15s %-15s %-20s\n",
			id, name, status, schedule, mode, created)
	}

	logger.Infof("Found %d backup jobs", len(jobs))
}

// selectAndDisplayJob shows interactive job selection and displays details
func selectAndDisplayJob(client *api.Client, jobs []api.BackupJob) error {
	// Create options for huh select
	var options []huh.Option[int]
	for _, job := range jobs {
		// Create a nice display string for each job
		status := "🟢"
		if job.Status == "paused" {
			status = "🟡"
		} else if job.Status != "active" {
			status = "🔴"
		}

		display := fmt.Sprintf("%s %s (%s) - %s", status, job.Name, job.Status, job.Schedule)
		options = append(options, huh.NewOption(display, job.ID))
	}

	// Add option to view all in table
	options = append(options, huh.NewOption("📋 View all jobs in table format", -1))

	var selectedJobID int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select a backup job to view details:").
				Description(fmt.Sprintf("Found %d backup jobs", len(jobs))).
				Options(options...).
				Value(&selectedJobID),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("job selection failed: %w", err)
	}

	// Handle table view option
	if selectedJobID == -1 {
		fmt.Println()
		displayJobsTable(jobs)
		return nil
	}

	// Get and display job details
	jobDetails, err := client.GetBackupJobDetails(selectedJobID)
	if err != nil {
		return fmt.Errorf("failed to get job details: %w", err)
	}

	displayJobDetails(jobDetails)
	return nil
}

// displayJobDetails shows detailed information about a backup job
func displayJobDetails(job *api.BackupJob) {
	fmt.Println()
	fmt.Printf("📋 Backup Job Details\n")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Printf("ID:                  %d\n", job.ID)
	fmt.Printf("Name:                %s\n", job.Name)
	fmt.Printf("Status:              %s\n", getStatusIcon(job.Status))
	fmt.Printf("Schedule:            %s\n", job.Schedule)
	fmt.Printf("Execution Mode:      %s\n", job.ExecutionMode)

	if job.ServerID != nil {
		fmt.Printf("Server ID:           %d\n", *job.ServerID)
	} else {
		fmt.Printf("Server ID:           Direct API\n")
	}

	fmt.Printf("Database Source ID:  %d\n", job.DatabaseSourceID)
	fmt.Printf("Storage Dest ID:     %d\n", job.StorageDestinationID)
	fmt.Printf("Retention Days:      %d\n", job.RetentionDays)
	fmt.Printf("Compression:         %s\n", job.Compression)
	fmt.Printf("Team ID:             %d\n", job.TeamID)
	fmt.Printf("Workspace ID:        %d\n", job.WorkspaceID)

	if job.Tags != nil && *job.Tags != "" {
		fmt.Printf("Tags:                %s\n", *job.Tags)
	}

	fmt.Printf("Created:             %s\n", job.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:             %s\n", job.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()
}

// getStatusIcon returns an icon for the job status
func getStatusIcon(status string) string {
	switch status {
	case "active":
		return "🟢 " + status
	case "paused":
		return "🟡 " + status
	case "failed":
		return "🔴 " + status
	default:
		return "⚪ " + status
	}
}

// truncateString truncates a string to the specified length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// formatBytes formats byte size in human readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
