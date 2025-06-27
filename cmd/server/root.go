package server

import (
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage servers",
	Long:  "Commands to manage and view servers",
}
