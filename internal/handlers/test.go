package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aawadall/mcpcli/internal/core"
)

// TestOptions contains flags for running tests.
type TestOptions struct {
	Config           string
	TestAll          bool
	TestResources    bool
	TestTools        bool
	TestCapabilities bool
	TestInit         bool
	ScriptFile       string
}

// LoadMCPConfig reads a configuration file which may be either an MCPConfig
// or a ProjectConfig and returns the resulting MCPConfig.
func LoadMCPConfig(configPath string) (*core.MCPConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var projectConfig core.ProjectConfig
	if err := json.Unmarshal(data, &projectConfig); err == nil {
		templateData := projectConfig.GetTemplateData()
		return templateData.MCPConfig, nil
	}
	var config core.MCPConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, core.FormatJSONError(data, err, "failed to parse config file")
	}
	return &config, nil
}

// RunTests connects to an MCP server based on the config and executes the
// selected tests.
func RunTests(opts *TestOptions, config *core.MCPConfig) error {
	var client *core.MCPClient
	if config.Transport.Type == "stdio" {
		if serverCmd, ok := config.Transport.Options["command"].(string); ok && serverCmd != "" {
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
			client = core.NewMCPClient()
		}
	} else {
		client = core.NewMCPClient()
	}

	id := 1
	if opts.ScriptFile != "" {
		fmt.Printf("⚠️ Reading and executing script: %s\n", opts.ScriptFile)
		// TODO: Read and execute script file
		return nil
	}

	if opts.TestAll || opts.TestResources {
		fmt.Printf("⚠️ Testing resources...\n")
		fmt.Printf("⚠️ Sending request: {\"method\":\"resources/list\",\"id\":%d}\n", id)
		resp, err := client.ListResources(id)
		id++
		if err != nil {
			fmt.Printf("❌ Failed to list resources: %v\n", err)
			fmt.Printf("⚠️ Make sure an MCP server is running and connected via stdin/stdout\n")
		} else if resp.Error != nil {
			fmt.Printf("❌ MCP error: %s\n", resp.Error.Message)
		} else {
			fmt.Printf("✅ Resources: %v\n", resp.Result)
		}
	}

	if opts.TestAll || opts.TestTools {
		fmt.Printf("⚠️ Testing tools...\n")
		fmt.Printf("⚠️ Sending request: {\"method\":\"tools/list\",\"id\":%d}\n", id)
		resp, err := client.ListTools(id)
		id++
		if err != nil {
			fmt.Printf("❌ Failed to list tools: %v\n", err)
			fmt.Printf("⚠️ Make sure an MCP server is running and connected via stdin/stdout\n")
		} else if resp.Error != nil {
			fmt.Printf("❌ MCP error: %s\n", resp.Error.Message)
		} else {
			fmt.Printf("✅ Tools: %v\n", resp.Result)
		}
	}

	// TODO: Add capabilities and init tests
	return nil
}
