package storage_destination

import (
	"github.com/spf13/cobra"
)

var StorageDestinationCmd = &cobra.Command{
	Use:   "storage-destination",
	Short: "Manage storage destinations",
	Long:  "Commands to manage and view storage destinations",
}