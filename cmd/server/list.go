package server

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
	Short: "List servers",
	Long:  "List all servers for the authenticated team. Use interactive selection to view server details.",
	Run:   runList,
}

func init() {
	listCmd.Flags().BoolVar(&tableView, "table", false, "Display servers in table format instead of interactive selection")
	ServerCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	logger.Debug("Starting servers list")

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		logger.Errorf("Failed to create API client: %v", err)
		return
	}

	// Get servers
	servers, err := client.GetServers()
	if err != nil {
		logger.Errorf("Failed to retrieve servers: %v", err)
		return
	}

	if len(servers) == 0 {
		logger.Info("No servers found")
		return
	}

	// Show table view if requested
	if tableView {
		displayServersTable(servers)
		return
	}

	// Interactive selection
	if err := selectAndDisplayServer(client, servers); err != nil {
		logger.Debugf("Interactive selection failed: %v", err)
		logger.Info("Falling back to table view...")
		fmt.Println()
		displayServersTable(servers)
	}
}

// displayServersTable shows servers in a table format
func displayServersTable(servers []api.Server) {
	fmt.Printf("%-5s %-20s %-15s %-10s %-15s %-20s\n", "ID", "NAME", "SSH HOST", "STATUS", "SSH USER", "CREATED")
	fmt.Println(strings.Repeat("-", 85))

	for _, server := range servers {
		id := fmt.Sprintf("%d", server.ID)
		name := truncateString(server.Name, 20)
		sshHost := "N/A"
		if server.SSHHost != nil {
			sshHost = truncateString(*server.SSHHost, 15)
		}
		status := truncateString(server.Status, 10)
		sshUser := truncateString(server.SSHUser, 15)
		created := server.CreatedAt.Format("2006-01-02 15:04")

		fmt.Printf("%-5s %-20s %-15s %-10s %-15s %-20s\n",
			id, name, sshHost, status, sshUser, created)
	}

	logger.Infof("Found %d servers", len(servers))
}

// selectAndDisplayServer shows interactive server selection and displays details
func selectAndDisplayServer(client *api.Client, servers []api.Server) error {
	// Create options for huh select
	var options []huh.Option[int]
	for _, server := range servers {
		// Create a nice display string for each server
		status := "🟢"
		if server.Status == "failing" {
			status = "🔴"
		} else if server.Status != "reachable" {
			status = "🟡"
		}

		sshInfo := "Direct API"
		if server.SSHHost != nil {
			sshInfo = fmt.Sprintf("%s@%s:%d", server.SSHUser, *server.SSHHost, server.SSHPort)
		}

		display := fmt.Sprintf("%s %s (%s) - %s", status, server.Name, server.Status, sshInfo)
		options = append(options, huh.NewOption(display, server.ID))
	}

	// Add option to view all in table
	options = append(options, huh.NewOption("📋 View all servers in table format", -1))

	var selectedServerID int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select a server to view details:").
				Description(fmt.Sprintf("Found %d servers", len(servers))).
				Options(options...).
				Value(&selectedServerID),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("server selection failed: %w", err)
	}

	// Handle table view option
	if selectedServerID == -1 {
		fmt.Println()
		displayServersTable(servers)
		return nil
	}

	// Get and display server details
	serverDetails, err := client.GetServerDetails(selectedServerID)
	if err != nil {
		return fmt.Errorf("failed to get server details: %w", err)
	}

	displayServerDetails(serverDetails)
	return nil
}

// displayServerDetails shows detailed information about a server
func displayServerDetails(server *api.Server) {
	fmt.Println()
	fmt.Printf("🖥️  Server Details\n")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Printf("ID:                  %d\n", server.ID)
	fmt.Printf("Name:                %s\n", server.Name)
	fmt.Printf("Status:              %s\n", getStatusIcon(server.Status))

	if server.SSHHost != nil {
		fmt.Printf("SSH Host:            %s\n", *server.SSHHost)
		fmt.Printf("SSH Port:            %d\n", server.SSHPort)
	} else {
		fmt.Printf("SSH Host:            Not configured\n")
	}

	fmt.Printf("SSH User:            %s\n", server.SSHUser)
	fmt.Printf("Working Directory:   %s\n", server.Workdir)
	fmt.Printf("Team ID:             %d\n", server.TeamID)
	fmt.Printf("Workspace ID:        %d\n", server.WorkspaceID)

	if server.LastSeenAt != nil {
		fmt.Printf("Last Seen:           %s\n", server.LastSeenAt.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Printf("Last Seen:           Never\n")
	}

	if server.Tags != nil && *server.Tags != "" {
		fmt.Printf("Tags:                %s\n", *server.Tags)
	}

	fmt.Printf("Created:             %s\n", server.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:             %s\n", server.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()
}

// getStatusIcon returns an icon for the server status
func getStatusIcon(status string) string {
	switch status {
	case "reachable":
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
