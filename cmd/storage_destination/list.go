package storage_destination

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
	Short: "List storage destinations",
	Long:  "List all storage destinations for the authenticated team. Use interactive selection to view destination details.",
	Run:   runList,
}

func init() {
	listCmd.Flags().BoolVar(&tableView, "table", false, "Display storage destinations in table format instead of interactive selection")
	StorageDestinationCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	logger.Debug("Starting storage destinations list")
	
	// Create API client
	client, err := api.NewClient()
	if err != nil {
		logger.Errorf("Failed to create API client: %v", err)
		return
	}
	
	// Get storage destinations
	destinations, err := client.GetStorageDestinations()
	if err != nil {
		logger.Errorf("Failed to retrieve storage destinations: %v", err)
		return
	}
	
	if len(destinations) == 0 {
		logger.Info("No storage destinations found")
		return
	}
	
	// Show table view if requested
	if tableView {
		displayDestinationsTable(destinations)
		return
	}
	
	// Interactive selection
	if err := selectAndDisplayDestination(client, destinations); err != nil {
		logger.Debugf("Interactive selection failed: %v", err)
		logger.Info("Falling back to table view...")
		fmt.Println()
		displayDestinationsTable(destinations)
	}
}

// displayDestinationsTable shows storage destinations in a table format
func displayDestinationsTable(destinations []api.StorageDestination) {
	fmt.Printf("%-5s %-10s %-15s %-25s %-20s %-9s %-15s\n", "ID", "DRIVER", "REGION", "BUCKET", "PREFIX", "DEFAULT", "CREATED")
	fmt.Println(strings.Repeat("-", 100))
	
	for _, dest := range destinations {
		id := fmt.Sprintf("%d", dest.ID)
		driver := truncateString(dest.Driver, 10)
		region := "N/A"
		if dest.Region != nil {
			region = truncateString(*dest.Region, 15)
		}
		bucket := truncateString(dest.Bucket, 25)
		prefix := "N/A"
		if dest.PathPrefix != nil {
			prefix = truncateString(*dest.PathPrefix, 20)
		}
		isDefault := "No"
		if dest.IsDefault {
			isDefault = "Yes"
		}
		created := dest.CreatedAt.Format("2006-01-02 15:04")
		
		fmt.Printf("%-5s %-10s %-15s %-25s %-20s %-9s %-15s\n", 
			id, driver, region, bucket, prefix, isDefault, created)
	}
	
	logger.Infof("Found %d storage destinations", len(destinations))
}

// selectAndDisplayDestination shows interactive storage destination selection and displays details
func selectAndDisplayDestination(client *api.Client, destinations []api.StorageDestination) error {
	// Create options for huh select
	var options []huh.Option[int]
	for _, dest := range destinations {
		// Create a nice display string for each destination
		status := "📦"
		if dest.IsDefault {
			status = "⭐"
		}
		
		bucketInfo := dest.Bucket
		if dest.Region != nil {
			bucketInfo += " (" + *dest.Region + ")"
		}
		
		display := fmt.Sprintf("%s %s - %s", status, dest.Driver, bucketInfo)
		options = append(options, huh.NewOption(display, dest.ID))
	}
	
	// Add option to view all in table
	options = append(options, huh.NewOption("📋 View all storage destinations in table format", -1))

	var selectedDestID int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select a storage destination to view details:").
				Description(fmt.Sprintf("Found %d storage destinations", len(destinations))).
				Options(options...).
				Value(&selectedDestID),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("storage destination selection failed: %w", err)
	}

	// Handle table view option
	if selectedDestID == -1 {
		fmt.Println()
		displayDestinationsTable(destinations)
		return nil
	}

	// Get and display destination details
	destDetails, err := client.GetStorageDestinationDetails(selectedDestID)
	if err != nil {
		return fmt.Errorf("failed to get storage destination details: %w", err)
	}

	displayDestinationDetails(destDetails)
	return nil
}

// displayDestinationDetails shows detailed information about a storage destination
func displayDestinationDetails(dest *api.StorageDestination) {
	fmt.Println()
	fmt.Printf("📦 Storage Destination Details\n")
	fmt.Println(strings.Repeat("=", 50))
	
	fmt.Printf("ID:                  %d\n", dest.ID)
	fmt.Printf("Driver:              %s\n", dest.Driver)
	fmt.Printf("Bucket:              %s\n", dest.Bucket)
	
	if dest.Region != nil {
		fmt.Printf("Region:              %s\n", *dest.Region)
	} else {
		fmt.Printf("Region:              Not specified\n")
	}
	
	if dest.PathPrefix != nil && *dest.PathPrefix != "" {
		fmt.Printf("Path Prefix:         %s\n", *dest.PathPrefix)
	} else {
		fmt.Printf("Path Prefix:         Not specified\n")
	}
	
	fmt.Printf("Is Default:          %t\n", dest.IsDefault)
	
	if dest.ConfigJSON != nil && *dest.ConfigJSON != "" {
		fmt.Printf("Configuration:       %s\n", "[CUSTOM CONFIG]")
	}
	
	fmt.Printf("Team ID:             %d\n", dest.TeamID)
	fmt.Printf("Workspace ID:        %d\n", dest.WorkspaceID)
	
	if dest.Tags != nil && *dest.Tags != "" {
		fmt.Printf("Tags:                %s\n", *dest.Tags)
	}
	
	fmt.Printf("Created:             %s\n", dest.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:             %s\n", dest.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()
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