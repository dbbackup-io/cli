package database_source

import (
	"github.com/spf13/cobra"
)

var DatabaseSourceCmd = &cobra.Command{
	Use:   "database-source",
	Short: "Manage database sources",
	Long:  "Commands to manage and view database sources",
}
