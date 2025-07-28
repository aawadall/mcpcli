package commands

import (
	"github.com/spf13/cobra"
)

// make root command
func MakeRootCommand(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "mcpcli",
		Short:   "A CLI tool for MCP (Model Context Protocol) development",
		Long:    `mcpcli is a command-line tool for scaffolding, testing, and managing MCP servers across multiple languages.`,
		Version: version,
	}

	// Add subcommands
	rootCmd.AddCommand(NewGenerateCmd())
	rootCmd.AddCommand(NewTestCmd())
	// TODO: Add future commands

	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress output")

	return rootCmd
}
