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

// GenerateOptions holds all configurable parameters for the generate command.
// When running in interactive mode, unanswered fields will be prompted.
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

// NewGenerateCmd creates the `generate` cobra command used to scaffold new MCP
// server projects. It sets up flags, validation, and interactive prompts.
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
	if err := askBasicOptions(opts); err != nil {
		return err
	}
	if err := promptForTools(opts); err != nil {
		return err
	}
	if err := promptForResources(opts); err != nil {
		return err
	}
	return promptForCapabilities(opts)
}

// askBasicOptions collects core project options interactively.
// It updates the provided GenerateOptions with any answers.
func askBasicOptions(opts *GenerateOptions) error {
	questions := buildBasicQuestions(opts)
	ans := basicAnswers{}
	if len(questions) > 0 {
		if err := survey.Ask(questions, &ans); err != nil {
			return err
		}
	}
	applyBasicAnswers(opts, ans)
	return nil
}

type basicAnswers struct {
	Name      string
	Language  string
	Transport string
	Docker    bool
	Examples  bool
	Output    string
}

// buildBasicQuestions returns survey questions for unset basic options.
func buildBasicQuestions(opts *GenerateOptions) []*survey.Question {
	qs := []*survey.Question{}
	if opts.Name == "" {
		qs = append(qs, &survey.Question{Name: "name", Prompt: &survey.Input{Message: "Project name:", Default: "my-mcp-server"}, Validate: survey.Required})
	}
	if opts.Language == "" {
		qs = append(qs, &survey.Question{Name: "language", Prompt: &survey.Select{Message: "Select programming language:", Options: []string{"golang", "python", "java", "javascript"}, Default: "golang"}, Validate: survey.Required})
	}
	if opts.Transport == "" {
		qs = append(qs, &survey.Question{Name: "transport", Prompt: &survey.Select{Message: "Choose transport method:", Options: []string{"stdio", "rest", "websocket"}, Default: "stdio"}})
	}
	if !opts.Docker {
		qs = append(qs, &survey.Question{Name: "docker", Prompt: &survey.Confirm{Message: "Include Docker support?", Default: true}})
	}
	if !opts.Examples {
		qs = append(qs, &survey.Question{Name: "examples", Prompt: &survey.Confirm{Message: "Include example resources and tools?", Default: true}})
	}
	if opts.Output == "" {
		qs = append(qs, &survey.Question{Name: "output", Prompt: &survey.Input{Message: "Output directory (leave empty to use project name):"}})
	}
	return qs
}

// applyBasicAnswers applies the survey answers back to the options struct.
func applyBasicAnswers(opts *GenerateOptions, ans basicAnswers) {
	if opts.Name == "" {
		opts.Name = ans.Name
	}
	if opts.Language == "" {
		opts.Language = ans.Language
	}
	if opts.Transport == "" {
		opts.Transport = ans.Transport
	}
	if !opts.Docker {
		opts.Docker = ans.Docker
	}
	if !opts.Examples {
		opts.Examples = ans.Examples
	}
	if opts.Output == "" {
		opts.Output = ans.Output
	}
}

// promptForTools interactively adds tool definitions to the options.
func promptForTools(opts *GenerateOptions) error {
	var add bool
	survey.AskOne(&survey.Confirm{Message: "Would you like to add tools?", Default: false}, &add)
	for add {
		var tool core.Tool
		survey.Ask([]*survey.Question{
			{Name: "Name", Prompt: &survey.Input{Message: "Tool name:"}, Validate: survey.Required},
			{Name: "Description", Prompt: &survey.Input{Message: "Tool description:"}},
		}, &tool)
		opts.Tools = append(opts.Tools, tool)
		survey.AskOne(&survey.Confirm{Message: "Add another tool?", Default: false}, &add)
	}
	return nil
}

// promptForResources interactively adds resource definitions to the options.
func promptForResources(opts *GenerateOptions) error {
	var add bool
	survey.AskOne(&survey.Confirm{Message: "Would you like to add resources?", Default: false}, &add)
	for add {
		var res core.Resource
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
		}, &res)
		opts.Resources = append(opts.Resources, res)
		survey.AskOne(&survey.Confirm{Message: "Add another resource?", Default: false}, &add)
	}
	return nil
}

// promptForCapabilities interactively adds capability definitions to the options.
func promptForCapabilities(opts *GenerateOptions) error {
	var add bool
	survey.AskOne(&survey.Confirm{Message: "Would you like to add capabilities?", Default: false}, &add)
	for add {
		var cap core.Capability
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
		}, &cap)
		opts.Capabilities = append(opts.Capabilities, cap)
		survey.AskOne(&survey.Confirm{Message: "Add another capability?", Default: false}, &add)
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

	if err := prepareDirectory(opts.Output, opts.Force); err != nil {
		return err
	}

	generator, err := selectGenerator(opts.Language)
	if err != nil {
		return err
	}
	if err := generator.Generate(config); err != nil {
		os.RemoveAll(opts.Output)
		return fmt.Errorf("failed to generate project: %w", err)
	}

	printNextSteps(opts)
	return nil
}

// prepareDirectory creates the output directory, removing it first when force is true.
func prepareDirectory(path string, force bool) error {
	if stat, err := os.Stat(path); err == nil {
		if force {
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("failed to remove existing directory: %w", err)
			}
		} else if stat != nil {
			return fmt.Errorf("output directory %s already exists, use --force to overwrite", path)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check output directory: %w", err)
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	return nil
}

// selectGenerator returns a project generator based on the chosen language.
func selectGenerator(lang string) (generators.Generator, error) {
	switch lang {
	case "golang":
		return generators.NewGolangGenerator(), nil
	case "javascript":
		return generators.NewNodeGenerator(), nil
	case "java":
		return generators.NewJavaGenerator(), nil
	case "python":
		return generators.NewPythonGenerator(), nil
	}
	return nil, fmt.Errorf("language %s is not supported yet", lang)
}

// printNextSteps outputs instructions for running the generated project.
func printNextSteps(opts *GenerateOptions) {
	fmt.Printf("✅ Successfully generated MCP server project: %s\n", opts.Name)
	path, _ := filepath.Abs(opts.Output)
	fmt.Printf("📁 Location: %s\n", path)
	fmt.Printf("🚀 Next steps:\n")
	fmt.Printf("   cd %s\n", opts.Output)

	switch opts.Language {
	case "go", "golang":
		fmt.Printf("   go mod tidy\n")
		fmt.Printf("   go run cmd/%s/main.go\n", opts.Transport)
	case "javascript":
		fmt.Printf("   npm install\n")
		fmt.Printf("   node src/index.js\n")
	case "java":
		fmt.Printf("   mvn package\n")
		fmt.Printf("   java -jar target/%s-1.0.0.jar\n", opts.Name)
	case "python":
		fmt.Printf("   python src/main.py\n")
	}
}
