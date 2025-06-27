package job

import (
	"github.com/spf13/cobra"
)

var JobCmd = &cobra.Command{
	Use:   "job",
	Short: "Manage backup jobs",
	Long:  "Commands to manage and view backup jobs",
}
