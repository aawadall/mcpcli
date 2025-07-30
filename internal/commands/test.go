// Package commands holds the Cobra commands used by mcpcli.
package commands

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aawadall/mcpcli/internal/handlers"
	"github.com/spf13/cobra"
)

// TestOptions contains flags for the `test` command.
type TestOptions = handlers.TestOptions

// needsTestInteractiveMode returns true if no test flags are set and
// the command should prompt the user interactively.
func needsTestInteractiveMode(opts *TestOptions) bool {
	return !opts.TestAll && !opts.TestResources && !opts.TestTools && !opts.TestCapabilities && !opts.TestInit && opts.ScriptFile == "" && opts.Config == ""
}

// promptForTestOptions displays an interactive survey to choose which tests to run.
func promptForTestOptions(opts *TestOptions) error {
	choices := []string{"Resources", "Tools", "Capabilities", "Initialization", "All"}
	selected := []string{}
	prompt := &survey.MultiSelect{
		Message: "Which tests would you like to run?",
		Options: choices,
	}
	if err := survey.AskOne(prompt, &selected); err != nil {
		return err
	}
	for _, choice := range selected {
		switch choice {
		case "All":
			opts.TestAll = true
		case "Resources":
			opts.TestResources = true
		case "Tools":
			opts.TestTools = true
		case "Capabilities":
			opts.TestCapabilities = true
		case "Initialization":
			opts.TestInit = true
		}
	}

	// Prompt for config file if not set
	if opts.Config == "" {
		configPrompt := &survey.Input{
			Message: "Enter path to MCP configuration file:",
			Default: "configs/mcp-config.json",
		}
		if err := survey.AskOne(configPrompt, &opts.Config); err != nil {
			return err
		}
	}

	return nil
}

// loadMCPConfig reads a configuration file which may be either an MCPConfig
// or a ProjectConfig and returns the resulting MCPConfig.

// NewTestCmd creates the `test` cobra command.
func NewTestCmd() *cobra.Command {
	opts := &TestOptions{}

	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test MCP server resources, tools, capabilities, and initialization.",
		Long:  `Test command connects to an MCP server and runs tests on resources, tools, capabilities, and initialization. Optionally, test instructions can be read from a script file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if needsTestInteractiveMode(opts) {
				if err := promptForTestOptions(opts); err != nil {
					return fmt.Errorf("interactive prompt failed: %w", err)
				}
			}

			config, err := handlers.LoadMCPConfig(opts.Config)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			return handlers.RunTests(opts, config)
		},
	}

	cmd.Flags().StringVarP(&opts.Config, "config", "c", "", "Path to MCP configuration file")
	cmd.Flags().BoolVar(&opts.TestAll, "all", false, "Test all components (resources, tools, capabilities, init)")
	cmd.Flags().BoolVar(&opts.TestResources, "resources", false, "Test resources")
	cmd.Flags().BoolVar(&opts.TestTools, "tools", false, "Test tools")
	cmd.Flags().BoolVar(&opts.TestCapabilities, "capabilities", false, "Test capabilities")
	cmd.Flags().BoolVar(&opts.TestInit, "init", false, "Test initialization")
	cmd.Flags().StringVarP(&opts.ScriptFile, "script", "f", "", "Path to test script file")

	return cmd
}
