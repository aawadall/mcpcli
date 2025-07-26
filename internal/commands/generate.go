package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aawadall/mcpcli/internal/core"
	"github.com/aawadall/mcpcli/internal/generators"
	"github.com/spf13/cobra"
)

type GenerateOptions struct {
	Name         string
	Language     string
	Transport    string
	Docker       bool
	Examples     bool
	Output       string
	Interactive  bool
	Force        bool
	Tools        []core.Tool
	Resources    []core.Resource
	Capabilities []core.Capability
}

func NewGenerateCmd() *cobra.Command {
	opts := &GenerateOptions{}

	cmd := &cobra.Command{
		Use:     "generate [name]",
		Aliases: []string{"gen", "g"},
		Short:   "Generate a new MCP server project",
		Long: `Generate scaffolds a new MCP server project with the specified configuration.
		Supports multiple languages, transport methods, and includes optional Docker support.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// If name provided as argument
			if len(args) > 0 {
				opts.Name = args[0]
			}

			// Check if we need interactive mode
			opts.Interactive = needsInteractiveMode(opts)

			if opts.Interactive {
				if err := promptForOptions(opts); err != nil {
					return fmt.Errorf("interactive prompt failed: %w", err)
				}
			}

			// Validate options
			if err := validateOptions(opts); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			// Generate the project
			return generateProject(opts)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "MCP project name")
	cmd.Flags().StringVarP(&opts.Language, "language", "l", "", "Programming language (e.g., golang, python, java)")
	cmd.Flags().StringVarP(&opts.Transport, "transport", "t", "", "Transport method (e.g., stdio, rest, websocket)")
	cmd.Flags().BoolVarP(&opts.Docker, "docker", "d", false, "Include Docker support")
	cmd.Flags().BoolVarP(&opts.Examples, "examples", "e", false, "Include example resources and tools")
	cmd.Flags().StringVarP(&opts.Output, "output", "o", "", "Output directory (default to project name)")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Overwrite existing directory")

	return cmd
}

// needsInteractiveMode checks if the options are incomplete and requires user input.
func needsInteractiveMode(opts *GenerateOptions) bool {
	return opts.Name == "" || opts.Language == "" || opts.Transport == ""
}

// promptForOptions prompts the user for missing options interactively.
func promptForOptions(opts *GenerateOptions) error {
	questions := []*survey.Question{}

	// Project name
	if opts.Name == "" {
		questions = append(questions, &survey.Question{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Project name:",
				Default: "my-mcp-server",
			},
			Validate: survey.Required,
		})
	}

	// Language selection
	if opts.Language == "" {
		questions = append(questions, &survey.Question{
			Name: "language",
			Prompt: &survey.Select{
				Message: "Select programming language:",
				Options: []string{"golang", "python", "java", "javascript"},
				Default: "golang",
			},
			Validate: survey.Required,
		})
	}

	// Transport method
	if opts.Transport == "" {
		questions = append(questions, &survey.Question{
			Name: "transport",
			Prompt: &survey.Select{
				Message: "Choose transport method:",
				Options: []string{"stdio", "rest", "websocket"},
				Default: "stdio",
			},
		})
	}

	// Docker support
	if !opts.Docker {
		questions = append(questions, &survey.Question{
			Name: "docker",
			Prompt: &survey.Confirm{
				Message: "Include Docker support?",
				Default: true,
			},
		})
	}

	// Examples
	if !opts.Examples {
		questions = append(questions, &survey.Question{
			Name: "examples",
			Prompt: &survey.Confirm{
				Message: "Include example resources and tools?",
				Default: true,
			},
		})
	}

	// Output directory
	if opts.Output == "" {
		questions = append(questions, &survey.Question{
			Name: "output",
			Prompt: &survey.Input{
				Message: "Output directory (leave empty to use project name):",
			},
		})
	}

	// Ask the questions
	answers := struct {
		Name      string
		Language  string
		Transport string
		Docker    bool
		Examples  bool
		Output    string
	}{}

	// Run the survey if there are questions
	if len(questions) > 0 {
		err := survey.Ask(questions, &answers)
		if err != nil {
			return err
		}
		// Assign answers to options
		if opts.Name == "" {
			opts.Name = answers.Name
		}

		if opts.Language == "" {
			opts.Language = answers.Language
		}

		if opts.Transport == "" {
			opts.Transport = answers.Transport
		}

		if !opts.Docker {
			opts.Docker = answers.Docker
		}

		if !opts.Examples {
			opts.Examples = answers.Examples
		}

		if opts.Output == "" {
			opts.Output = answers.Output
		}
	}

	var addTools bool
	survey.AskOne(&survey.Confirm{
		Message: "Would you like to add tools?",
		Default: false,
	}, &addTools)

	for addTools {
		var tool core.Tool
		survey.Ask([]*survey.Question{
			{
				Name:     "Name",
				Prompt:   &survey.Input{Message: "Tool name:"},
				Validate: survey.Required,
			},
			{
				Name:   "Description",
				Prompt: &survey.Input{Message: "Tool description:"},
			},
		}, &tool)
		opts.Tools = append(opts.Tools, tool)

		survey.AskOne(&survey.Confirm{
			Message: "Add another tool?",
			Default: false,
		}, &addTools)
	}

	var addResources bool
	survey.AskOne(&survey.Confirm{
		Message: "Would you like to add resources?",
		Default: false,
	}, &addResources)

	for addResources {
		var resource core.Resource
		survey.Ask([]*survey.Question{
			{
				Name:     "Name",
				Prompt:   &survey.Input{Message: "Resource name:"},
				Validate: survey.Required,
			},
			{
				Name: "Type",
				Prompt: &survey.Select{
					Message: "Resource type:",
					Options: []string{string(core.ResourceTypeDatabase), string(core.ResourceTypeFilesystem), string(core.ResourceTypeTime)},
				},
			},
		}, &resource)
		opts.Resources = append(opts.Resources, resource)

		survey.AskOne(&survey.Confirm{
			Message: "Add another resource?",
			Default: false,
		}, &addResources)
	}

	var addCapabilities bool
	survey.AskOne(&survey.Confirm{
		Message: "Would you like to add capabilities?",
		Default: false,
	}, &addCapabilities)

	for addCapabilities {
		var capability core.Capability
		survey.Ask([]*survey.Question{
			{
				Name:     "Name",
				Prompt:   &survey.Input{Message: "Capability name:"},
				Validate: survey.Required,
			},
			{
				Name:   "Enabled",
				Prompt: &survey.Confirm{Message: "Enable this capability?"},
			},
		}, &capability)
		opts.Capabilities = append(opts.Capabilities, capability)

		survey.AskOne(&survey.Confirm{
			Message: "Add another capability?",
			Default: false,
		}, &addCapabilities)
	}

	return nil
}

// validateOptions checks if the provided options are valid.
func validateOptions(opts *GenerateOptions) error {
	// validate name
	if opts.Name == "" {
		return fmt.Errorf("project name is required")
	}

	// validate language
	validLanguages := []string{"golang", "javascript", "java", "python"}
	if !contains(validLanguages, opts.Language) {
		return fmt.Errorf("invalid language: %s, valid options are: %v", opts.Language, validLanguages)
	}

	// validate transport
	validTransports := []string{"stdio"} // Add other valid transports as needed
	if !contains(validTransports, opts.Transport) {
		return fmt.Errorf("invalid transport: %s, valid options are: %v", opts.Transport, validTransports)
	}

	// Set default output directory if not provided
	if opts.Output == "" {
		opts.Output = opts.Name
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// generateProject creates the project structure based on the provided options.
func generateProject(opts *GenerateOptions) error {
	// Create project configuration
	config := &core.ProjectConfig{
		Name:         opts.Name,
		Language:     opts.Language,
		Transport:    opts.Transport,
		Docker:       opts.Docker,
		Examples:     opts.Examples,
		Output:       opts.Output,
		Tools:        opts.Tools,
		Resources:    opts.Resources,
		Capabilities: opts.Capabilities,
	}

	// Check if the directory exists and handle the --force flag
	if _, err := os.Stat(opts.Output); err == nil {
		if opts.Force {
			if err := os.RemoveAll(opts.Output); err != nil {
				return fmt.Errorf("failed to remove existing directory: %w", err)
			}
		} else {
			return fmt.Errorf("output directory %s already exists, use --force to overwrite", opts.Output)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check output directory: %w", err)
	}

	// Create output directory
	if err := os.MkdirAll(opts.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate project files based on the language
	var generator generators.Generator

	switch opts.Language {
	case "golang":
		generator = generators.NewGolangGenerator()
	case "javascript":
		generator = generators.NewNodeGenerator()
	case "java":
		generator = generators.NewJavaGenerator()
	case "python":
		generator = generators.NewPythonGenerator()
	default:
		return fmt.Errorf("language %s is not supported yet", opts.Language)
	}

	// Generate the project files
	if err := generator.Generate(config); err != nil {
		// rollback the created directory if generation fails
		os.RemoveAll(opts.Output)
		return fmt.Errorf("failed to generate project: %w", err)
	}

	// Success message
	fmt.Printf("✅ Successfully generated MCP server project: %s\n", opts.Name)
	path, _ := filepath.Abs(opts.Output)
	fmt.Printf("📁 Location: %s\n", path)
	fmt.Printf("🚀 Next steps:\n")
	fmt.Printf("   cd %s\n", opts.Output)

	if opts.Language == "go" || opts.Language == "golang" {
		fmt.Printf("   go mod tidy\n")
		fmt.Printf("   go run cmd/%s/main.go\n", opts.Transport)
	} else if opts.Language == "javascript" {
		fmt.Printf("   npm install\n")
		fmt.Printf("   node src/index.js\n")
	} else if opts.Language == "java" {
		fmt.Printf("   mvn package\n")
		fmt.Printf("   java -jar target/%s-1.0.0.jar\n", opts.Name)
	} else if opts.Language == "python" {
		fmt.Printf("   python src/main.py\n")
	}

	return nil
}
