package job

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/dbbackup-io/cli/pkg/api"
	"github.com/dbbackup-io/cli/pkg/logger"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [job-id]",
	Short: "Run a backup job",
	Long:  "Trigger a backup job to run immediately. You can provide a job ID or select interactively.",
	Args:  cobra.MaximumNArgs(1),
	Run:   runJob,
}

func init() {
	JobCmd.AddCommand(runCmd)
}

func runJob(cmd *cobra.Command, args []string) {
	logger.Debug("Starting job run")

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		logger.Errorf("Failed to create API client: %v", err)
		return
	}

	var jobID int

	// If job ID provided as argument, use it
	if len(args) > 0 {
		parsedID, err := strconv.Atoi(args[0])
		if err != nil {
			logger.Errorf("Invalid job ID: %v", err)
			return
		}
		jobID = parsedID
	} else {
		// Otherwise, let user select from available jobs
		jobs, err := client.GetBackupJobs()
		if err != nil {
			logger.Errorf("Failed to retrieve backup jobs: %v", err)
			return
		}

		if len(jobs) == 0 {
			logger.Info("No backup jobs found")
			return
		}

		selectedJob, err := selectJobToRun(jobs)
		if err != nil {
			logger.Errorf("Failed to select job: %v", err)
			return
		}

		jobID = selectedJob.ID
	}

	logger.Infof("Triggering backup job %d...", jobID)

	// Run the backup job
	if err := client.RunBackupJob(jobID); err != nil {
		logger.Errorf("Failed to run backup job: %v", err)
		return
	}

	logger.Infof("✅ Backup job %d has been triggered successfully", jobID)
}

// selectJobToRun shows interactive job selection for running
func selectJobToRun(jobs []api.BackupJob) (*api.BackupJob, error) {
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
				Title("Select a backup job to run:").
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
