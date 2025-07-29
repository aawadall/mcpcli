package main

import (
	"fmt"
	"os"

	"github.com/aawadall/mcpcli/internal/commands"
)

// version is the current release tag for the CLI.
var version = "0.4.2"

func main() {
	// Create the root command
	rootCmd := commands.MakeRootCommand(version)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
