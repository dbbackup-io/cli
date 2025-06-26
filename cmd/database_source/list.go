package database_source

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
	Short: "List database sources",
	Long:  "List all database sources for the authenticated team. Use interactive selection to view source details.",
	Run:   runList,
}

func init() {
	listCmd.Flags().BoolVar(&tableView, "table", false, "Display database sources in table format instead of interactive selection")
	DatabaseSourceCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	logger.Debug("Starting database sources list")
	
	// Create API client
	client, err := api.NewClient()
	if err != nil {
		logger.Errorf("Failed to create API client: %v", err)
		return
	}
	
	// Get database sources
	sources, err := client.GetDatabaseSources()
	if err != nil {
		logger.Errorf("Failed to retrieve database sources: %v", err)
		return
	}
	
	if len(sources) == 0 {
		logger.Info("No database sources found")
		return
	}
	
	// Show table view if requested
	if tableView {
		displaySourcesTable(sources)
		return
	}
	
	// Interactive selection
	if err := selectAndDisplaySource(client, sources); err != nil {
		logger.Debugf("Interactive selection failed: %v", err)
		logger.Info("Falling back to table view...")
		fmt.Println()
		displaySourcesTable(sources)
	}
}

// displaySourcesTable shows database sources in a table format
func displaySourcesTable(sources []api.DatabaseSource) {
	fmt.Printf("%-5s %-20s %-12s %-20s %-6s %-10s %-15s\n", "ID", "ENGINE", "HOST", "DATABASE", "PORT", "STATUS", "CREATED")
	fmt.Println(strings.Repeat("-", 88))
	
	for _, source := range sources {
		id := fmt.Sprintf("%d", source.ID)
		engine := truncateString(source.Engine, 12)
		host := truncateString(source.Host, 20)
		dbName := "N/A"
		if source.DBName != nil {
			dbName = truncateString(*source.DBName, 20)
		}
		port := fmt.Sprintf("%d", source.Port)
		status := truncateString(source.Status, 10)
		created := source.CreatedAt.Format("2006-01-02 15:04")
		
		fmt.Printf("%-5s %-20s %-12s %-20s %-6s %-10s %-15s\n", 
			id, engine, host, dbName, port, status, created)
	}
	
	logger.Infof("Found %d database sources", len(sources))
}

// selectAndDisplaySource shows interactive database source selection and displays details
func selectAndDisplaySource(client *api.Client, sources []api.DatabaseSource) error {
	// Create options for huh select
	var options []huh.Option[int]
	for _, source := range sources {
		// Create a nice display string for each source
		status := "🟢"
		if source.Status == "failing" {
			status = "🔴"
		} else if source.Status != "active" {
			status = "🟡"
		}
		
		dbInfo := fmt.Sprintf("%s:%d", source.Host, source.Port)
		if source.DBName != nil {
			dbInfo += "/" + *source.DBName
		}
		
		display := fmt.Sprintf("%s %s (%s) - %s", status, source.Engine, source.Status, dbInfo)
		options = append(options, huh.NewOption(display, source.ID))
	}
	
	// Add option to view all in table
	options = append(options, huh.NewOption("📋 View all database sources in table format", -1))

	var selectedSourceID int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select a database source to view details:").
				Description(fmt.Sprintf("Found %d database sources", len(sources))).
				Options(options...).
				Value(&selectedSourceID),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("database source selection failed: %w", err)
	}

	// Handle table view option
	if selectedSourceID == -1 {
		fmt.Println()
		displaySourcesTable(sources)
		return nil
	}

	// Get and display source details
	sourceDetails, err := client.GetDatabaseSourceDetails(selectedSourceID)
	if err != nil {
		return fmt.Errorf("failed to get database source details: %w", err)
	}

	displaySourceDetails(sourceDetails)
	return nil
}

// displaySourceDetails shows detailed information about a database source
func displaySourceDetails(source *api.DatabaseSource) {
	fmt.Println()
	fmt.Printf("🗄️  Database Source Details\n")
	fmt.Println(strings.Repeat("=", 50))
	
	fmt.Printf("ID:                  %d\n", source.ID)
	fmt.Printf("Engine:              %s\n", source.Engine)
	fmt.Printf("Host:                %s\n", source.Host)
	fmt.Printf("Port:                %d\n", source.Port)
	fmt.Printf("Status:              %s\n", getStatusIcon(source.Status))
	
	if source.DBName != nil {
		fmt.Printf("Database Name:       %s\n", *source.DBName)
	} else {
		fmt.Printf("Database Name:       Not specified\n")
	}
	
	if source.Username != nil {
		fmt.Printf("Username:            %s\n", *source.Username)
	} else {
		fmt.Printf("Username:            Not configured\n")
	}
	
	if source.Password != nil {
		fmt.Printf("Password:            %s\n", "[CONFIGURED]")
	} else {
		fmt.Printf("Password:            Not configured\n")
	}
	
	if source.ConfigJSON != nil && *source.ConfigJSON != "" {
		fmt.Printf("Configuration:       %s\n", "[CUSTOM CONFIG]")
	}
	
	fmt.Printf("Team ID:             %d\n", source.TeamID)
	fmt.Printf("Workspace ID:        %d\n", source.WorkspaceID)
	
	if source.Tags != nil && *source.Tags != "" {
		fmt.Printf("Tags:                %s\n", *source.Tags)
	}
	
	fmt.Printf("Created:             %s\n", source.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:             %s\n", source.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()
}

// getStatusIcon returns an icon for the database source status
func getStatusIcon(status string) string {
	switch status {
	case "active":
		return "🟢 " + status
	case "failing":
		return "🔴 " + status
	default:
		return "🟡 " + status
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