package main

import (
	"fmt"
	"os"

	"github.com/aawadall/mcpcli/internal/commands"
)

// version is the current release tag for the CLI.
var version = "0.4.1"

func main() {
	// Create the root command
	rootCmd := commands.MakeRootCommand(version)
	rootCmd.SetArgs(os.Args[1:])

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
