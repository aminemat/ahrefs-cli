package main

import (
	"os"

	"github.com/amine/ahrefs-cli/cmd"
	"github.com/amine/ahrefs-cli/cmd/config"
	"github.com/amine/ahrefs-cli/cmd/siteexplorer"
)

func main() {
	// Register all subcommands
	cmd.AddCommands(
		config.NewConfigCmd(),
		siteexplorer.NewSiteExplorerCmd(),
	)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
