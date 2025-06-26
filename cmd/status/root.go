package status

import (
	"fmt"

	"github.com/dbbackup-io/cli/pkg/config"
	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  "Display current authentication status including token and team information",
	Run:   runStatus,
}

func runStatus(cmd *cobra.Command, args []string) {
	if !config.IsAuthenticated() {
		fmt.Println("Status: Not authenticated")
		return
	}

	token := config.GetToken()
	teamId := config.GetTeamId()

	fmt.Println("Status: Authenticated")
	if token != "" {
		// Only show first 10 chars of token for security
		if len(token) > 10 {
			fmt.Printf("Token: %s...\n", token[:10])
		} else {
			fmt.Printf("Token: %s\n", token)
		}
	}
	if teamId != "" {
		fmt.Printf("Team ID: %s\n", teamId)
	}
}