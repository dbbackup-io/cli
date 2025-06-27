package main

import (
	"github.com/dbbackup-io/cli/cmd"
)

// version is set during build time
var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
