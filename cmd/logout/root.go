package logout

import (
	"fmt"

	"github.com/dbbackup-io/cli/pkg/config"
	"github.com/spf13/cobra"
)

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Sign out from DBBackup",
	Long:  "Remove the authentication token and sign out from DBBackup",
	Run:   runLogout,
}

func runLogout(cmd *cobra.Command, args []string) {
	if !config.IsAuthenticated() {
		fmt.Println("Not currently authenticated.")
		return
	}

	if err := config.ClearAuthData(); err != nil {
		fmt.Printf("Failed to clear authentication data: %v\n", err)
		return
	}

	fmt.Println("Successfully signed out!")
}