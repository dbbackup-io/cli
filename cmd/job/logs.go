package job

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/dbbackup-io/cli/pkg/api"
	"github.com/dbbackup-io/cli/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	logsTableView bool
)

var logsCmd = &cobra.Command{
	Use:   "logs [job-id]",
	Short: "View backup runs for a job",
	Long:  "List backup runs for a job. If no job ID is provided, you'll be prompted to select one interactively.",
	Args:  cobra.MaximumNArgs(1),
	Run:   runLogs,
}

func init() {
	logsCmd.Flags().BoolVar(&logsTableView, "table", false, "Display backup runs in table format instead of interactive selection")
	JobCmd.AddCommand(logsCmd)
}

func runLogs(cmd *cobra.Command, args []string) {
	logger.Debug("Starting job logs")

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		logger.Errorf("Failed to create API client: %v", err)
		return
	}

	var jobID int
	var selectedJob *api.BackupJob

	// If job ID provided as argument, use it
	if len(args) > 0 {
		parsedID, err := strconv.Atoi(args[0])
		if err != nil {
			logger.Errorf("Invalid job ID: %v", err)
			return
		}
		jobID = parsedID

		// Get job details for display
		jobDetails, err := client.GetBackupJobDetails(jobID)
		if err != nil {
			logger.Errorf("Failed to get job details: %v", err)
			return
		}
		selectedJob = jobDetails
	} else {
		// Let user select from available jobs
		jobs, err := client.GetBackupJobs()
		if err != nil {
			logger.Errorf("Failed to retrieve backup jobs: %v", err)
			return
		}

		if len(jobs) == 0 {
			logger.Info("No backup jobs found")
			return
		}

		selected, err := selectJobForLogs(jobs)
		if err != nil {
			logger.Errorf("Failed to select job: %v", err)
			return
		}

		jobID = selected.ID
		selectedJob = selected
	}

	// Get backup runs for the selected job
	jobRuns, err := client.GetBackupRunsForJob(jobID)
	if err != nil {
		logger.Errorf("Failed to retrieve backup runs for job %d: %v", jobID, err)
		return
	}

	if len(jobRuns) == 0 {
		logger.Infof("No backup runs found for job %d (%s)", jobID, selectedJob.Name)
		return
	}

	logger.Infof("Found %d backup runs for job: %s", len(jobRuns), selectedJob.Name)

	// Show table view if requested
	if logsTableView {
		displayJobRunsTable(jobRuns, selectedJob)
		return
	}

	// Interactive selection
	if err := selectAndDisplayJobRun(client, jobRuns, selectedJob); err != nil {
		logger.Debugf("Interactive selection failed: %v", err)
		logger.Info("Falling back to table view...")
		fmt.Println()
		displayJobRunsTable(jobRuns, selectedJob)
	}
}

// selectJobForLogs shows interactive job selection for viewing logs
func selectJobForLogs(jobs []api.BackupJob) (*api.BackupJob, error) {
	// Create options for huh select
	var options []huh.Option[int]
	for _, job := range jobs {
		// Create a nice display string for each job
		status := "🟢"
		if job.Status == "inactive" {
			status = "🔴"
		} else if job.Status != "active" {
			status = "🟡"
		}

		display := fmt.Sprintf("%s %s (%s) - Schedule: %s", status, job.Name, job.Status, job.Schedule)
		options = append(options, huh.NewOption(display, job.ID))
	}

	var selectedJobID int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select a backup job to view runs:").
				Description(fmt.Sprintf("Found %d backup jobs", len(jobs))).
				Options(options...).
				Value(&selectedJobID),
		),
	)

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("job selection failed: %w", err)
	}

	// Find the selected job
	for _, job := range jobs {
		if job.ID == selectedJobID {
			return &job, nil
		}
	}

	return nil, fmt.Errorf("selected job not found")
}

// displayJobRunsTable shows backup runs in a table format
func displayJobRunsTable(runs []api.BackupRun, job *api.BackupJob) {
	fmt.Printf("\n📋 Backup Runs for Job: %s (ID: %d)\n", job.Name, job.ID)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("%-5s %-12s %-20s %-20s %-10s %-8s\n", "ID", "STATUS", "STARTED", "FINISHED", "SIZE", "RETRIES")
	fmt.Println(strings.Repeat("-", 85))

	for _, run := range runs {
		id := fmt.Sprintf("%d", run.ID)
		status := truncateString(run.Status, 12)
		started := run.StartedAt.Format("2006-01-02 15:04")
		finished := "Running..."
		if run.FinishedAt != nil {
			finished = run.FinishedAt.Format("2006-01-02 15:04")
		}
		size := "N/A"
		if run.SizeBytes != nil {
			size = formatBytes(*run.SizeBytes)
		}
		retries := fmt.Sprintf("%d", run.RetryCount)

		fmt.Printf("%-5s %-12s %-20s %-20s %-10s %-8s\n",
			id, status, started, finished, size, retries)
	}

	fmt.Printf("\nTotal: %d backup runs\n", len(runs))
}

// selectAndDisplayJobRun shows interactive backup run selection and displays details
func selectAndDisplayJobRun(client *api.Client, runs []api.BackupRun, job *api.BackupJob) error {
	// Create options for huh select
	var options []huh.Option[int]
	for _, run := range runs {
		// Create a nice display string for each run
		status := getJobRunStatusIcon(run.Status)

		duration := "N/A"
		if run.RuntimeSec != nil {
			duration = fmt.Sprintf("%ds", *run.RuntimeSec)
		}

		size := "N/A"
		if run.SizeBytes != nil {
			size = formatBytes(*run.SizeBytes)
		}

		display := fmt.Sprintf("%s Run %d - %s (%s, %s)", status, run.ID, run.Status, duration, size)
		options = append(options, huh.NewOption(display, run.ID))
	}

	// Add option to view all in table
	options = append(options, huh.NewOption("📋 View all runs in table format", -1))

	var selectedRunID int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title(fmt.Sprintf("Select a backup run for job '%s':", job.Name)).
				Description(fmt.Sprintf("Found %d backup runs", len(runs))).
				Options(options...).
				Value(&selectedRunID),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("backup run selection failed: %w", err)
	}

	// Handle table view option
	if selectedRunID == -1 {
		fmt.Println()
		displayJobRunsTable(runs, job)
		return nil
	}

	// Get and display run details
	runDetails, err := client.GetBackupRunDetails(selectedRunID)
	if err != nil {
		return fmt.Errorf("failed to get backup run details: %w", err)
	}

	displayJobRunDetails(runDetails, job)
	return nil
}

// displayJobRunDetails shows detailed information about a backup run
func displayJobRunDetails(run *api.BackupRun, job *api.BackupJob) {
	fmt.Println()
	fmt.Printf("🔄 Backup Run Details\n")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Printf("Run ID:              %d\n", run.ID)
	fmt.Printf("Job:                 %s (ID: %d)\n", job.Name, job.ID)
	fmt.Printf("Status:              %s\n", getJobRunStatusIcon(run.Status))
	fmt.Printf("Started At:          %s\n", run.StartedAt.Format("2006-01-02 15:04:05"))

	if run.FinishedAt != nil {
		fmt.Printf("Finished At:         %s\n", run.FinishedAt.Format("2006-01-02 15:04:05"))
		duration := run.FinishedAt.Sub(run.StartedAt)
		fmt.Printf("Duration:            %s\n", duration.Round(time.Second))
	} else {
		fmt.Printf("Finished At:         Still running...\n")
		duration := time.Since(run.StartedAt)
		fmt.Printf("Running For:         %s\n", duration.Round(time.Second))
	}

	if run.RuntimeSec != nil {
		fmt.Printf("Runtime (seconds):   %d\n", *run.RuntimeSec)
	}

	if run.SizeBytes != nil {
		fmt.Printf("Backup Size:         %s (%d bytes)\n", formatBytes(*run.SizeBytes), *run.SizeBytes)
	} else {
		fmt.Printf("Backup Size:         Not available\n")
	}

	fmt.Printf("Retry Count:         %d\n", run.RetryCount)

	if run.LogPath != nil {
		fmt.Printf("Log Path:            %s\n", *run.LogPath)
	} else {
		fmt.Printf("Log Path:            Not available\n")
	}

	if run.BackupPath != nil {
		fmt.Printf("Backup Path:         %s\n", *run.BackupPath)
	} else {
		fmt.Printf("Backup Path:         Not available\n")
	}

	fmt.Printf("Created:             %s\n", run.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:             %s\n", run.UpdatedAt.Format("2006-01-02 15:04:05"))

	// Show job configuration
	fmt.Println()
	fmt.Printf("📝 Job Configuration\n")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Printf("Schedule:            %s\n", job.Schedule)
	fmt.Printf("Execution Mode:      %s\n", job.ExecutionMode)
	fmt.Printf("Retention Days:      %d\n", job.RetentionDays)
	fmt.Printf("Compression:         %s\n", job.Compression)

	fmt.Println()
}

// getJobRunStatusIcon returns an icon for the backup run status
func getJobRunStatusIcon(status string) string {
	switch status {
	case "completed", "success":
		return "✅ " + status
	case "failed", "error":
		return "❌ " + status
	case "running", "in_progress":
		return "🔄 " + status
	case "queued", "pending":
		return "⏳ " + status
	default:
		return "⚪ " + status
	}
}
