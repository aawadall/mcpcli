package main

import (
	"fmt"
	"os"

	"github.com/aawadall/mcpcli/internal/commands"
	"github.com/spf13/cobra"
)

// version is the current release tag for the CLI.
var version = "0.4.1"

func main() {
	rootCmd := &cobra.Command{
		Use:     "mcpcli",
		Short:   "A CLI tool for MCP (Model Context Protocol) development",
		Long:    `mcpcli is a command-line tool for scaffolding, testing, and managing MCP servers across multiple languages.`,
		Version: version,
	}

	// Add subcommands
	rootCmd.AddCommand(commands.NewGenerateCmd())
	rootCmd.AddCommand(commands.NewTestCmd())
	// TODO: Add future commands

	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress output")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
