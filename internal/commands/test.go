package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aawadall/mcpcli/internal/core"
	"github.com/spf13/cobra"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	emojiCheck  = "✅"
	emojiCross  = "❌"
	emojiWarn   = "⚠️"
)

type TestOptions struct {
	Config           string
	TestAll          bool
	TestResources    bool
	TestTools        bool
	TestCapabilities bool
	TestInit         bool
	ScriptFile       string
}

func needsTestInteractiveMode(opts *TestOptions) bool {
	return !opts.TestAll && !opts.TestResources && !opts.TestTools && !opts.TestCapabilities && !opts.TestInit && opts.ScriptFile == "" && opts.Config == ""
}

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

func loadMCPConfig(configPath string) (*core.MCPConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config core.MCPConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

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

			// Load MCP configuration
			config, err := loadMCPConfig(opts.Config)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			var client *core.MCPClient
			if config.Transport.Type == "stdio" {
				// For stdio, check if a command is specified to run the server
				if serverCmd, ok := config.Transport.Options["command"].(string); ok && serverCmd != "" {
					// Run the specified server command
					parts := strings.Fields(serverCmd)
					if len(parts) == 0 {
						return fmt.Errorf("invalid server command: %s", serverCmd)
					}
					cmd := exec.Command(parts[0], parts[1:]...)
					stdin, err := cmd.StdinPipe()
					if err != nil {
						return fmt.Errorf("failed to create stdin pipe: %w", err)
					}
					stdout, err := cmd.StdoutPipe()
					if err != nil {
						return fmt.Errorf("failed to create stdout pipe: %w", err)
					}
					if err := cmd.Start(); err != nil {
						return fmt.Errorf("failed to start server: %w", err)
					}
					defer cmd.Wait()
					client = core.NewMCPClientWithIO(stdout, stdin, os.Stderr)
				} else {
					// No command specified, assume server is already running via stdin/stdout
					client = core.NewMCPClient()
				}
			} else {
				client = core.NewMCPClient()
			}

			id := 1

			if opts.ScriptFile != "" {
				fmt.Printf("%s %sReading and executing script: %s%s\n", colorYellow, emojiWarn, opts.ScriptFile, colorReset)
				// TODO: Read and execute script file
				return nil
			}

			if opts.TestAll || opts.TestResources {
				fmt.Printf("%s Testing resources...%s\n", colorYellow, colorReset)
				fmt.Printf("%s Sending request: {\"method\":\"resources/list\",\"id\":%d}%s\n", colorYellow, id, colorReset)
				resp, err := client.ListResources(id)
				id++
				if err != nil {
					fmt.Printf("%s %sFailed to list resources: %v%s\n", colorRed, emojiCross, err, colorReset)
					fmt.Printf("%s %sMake sure an MCP server is running and connected via stdin/stdout%s\n", colorYellow, emojiWarn, colorReset)
				} else if resp.Error != nil {
					fmt.Printf("%s %sMCP error: %s%s\n", colorRed, emojiCross, resp.Error.Message, colorReset)
				} else {
					fmt.Printf("%s %sResources: %v%s\n", colorGreen, emojiCheck, resp.Result, colorReset)
				}
			}

			if opts.TestAll || opts.TestTools {
				fmt.Printf("%s Testing tools...%s\n", colorYellow, colorReset)
				fmt.Printf("%s Sending request: {\"method\":\"tools/list\",\"id\":%d}%s\n", colorYellow, id, colorReset)
				resp, err := client.ListTools(id)
				id++
				if err != nil {
					fmt.Printf("%s %sFailed to list tools: %v%s\n", colorRed, emojiCross, err, colorReset)
					fmt.Printf("%s %sMake sure an MCP server is running and connected via stdin/stdout%s\n", colorYellow, emojiWarn, colorReset)
				} else if resp.Error != nil {
					fmt.Printf("%s %sMCP error: %s%s\n", colorRed, emojiCross, resp.Error.Message, colorReset)
				} else {
					fmt.Printf("%s %sTools: %v%s\n", colorGreen, emojiCheck, resp.Result, colorReset)
				}
			}

			// TODO: Add capabilities and init tests

			return nil
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
