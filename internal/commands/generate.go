package commands

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aawadall/mcpcli/internal/core"
	"github.com/aawadall/mcpcli/internal/handlers"
	"github.com/spf13/cobra"
)

// NewGenerateCmd creates the `generate` cobra command used to scaffold new MCP
// server projects. It sets up flags, validation, and interactive prompts.
func NewGenerateCmd() *cobra.Command {
	opts := &handlers.GenerateOptions{}

	cmd := &cobra.Command{
		Use:     "generate [name]",
		Aliases: []string{"gen", "g"},
		Short:   "Generate a new MCP server project",
		Long: `Generate scaffolds a new MCP server project with the specified configuration.
Supports multiple languages, transport methods, and includes optional Docker support.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Name = args[0]
			}
			opts.Interactive = needsInteractiveMode(opts)
			if opts.Interactive {
				if err := promptForOptions(opts); err != nil {
					return fmt.Errorf("interactive prompt failed: %w", err)
				}
			}
			if err := handlers.ValidateGenerateOptions(opts); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}
			return handlers.GenerateProject(opts)
		},
	}

	addFlags(cmd, opts)
	return cmd
}

// Add flags to the generate command.
func addFlags(cmd *cobra.Command, opts *handlers.GenerateOptions) {
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "MCP project name")
	cmd.Flags().StringVarP(&opts.Language, "language", "l", "", "Programming language (e.g., golang, python, java)")
	cmd.Flags().StringVarP(&opts.Transport, "transport", "t", "", "Transport method (e.g., stdio, rest, websocket)")
	cmd.Flags().BoolVarP(&opts.Docker, "docker", "d", false, "Include Docker support")
	cmd.Flags().BoolVarP(&opts.Examples, "examples", "e", false, "Include example resources and tools")
	cmd.Flags().StringVarP(&opts.Output, "output", "o", "", "Output directory (default to project name)")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Overwrite existing directory")
}

// needsInteractiveMode checks if the options are incomplete and requires user input.
func needsInteractiveMode(opts *handlers.GenerateOptions) bool {
	return opts.Name == "" || opts.Language == "" || opts.Transport == ""
}

// promptForOptions prompts the user for missing options interactively.
func promptForOptions(opts *handlers.GenerateOptions) error {
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

type basicAnswers struct {
	Name      string
	Language  string
	Transport string
	Docker    bool
	Examples  bool
	Output    string
}

// askBasicOptions collects core project options interactively.
func askBasicOptions(opts *handlers.GenerateOptions) error {
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

// buildBasicQuestions returns survey questions for unset basic options.
func buildBasicQuestions(opts *handlers.GenerateOptions) []*survey.Question {
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
func applyBasicAnswers(opts *handlers.GenerateOptions, ans basicAnswers) {
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
func promptForTools(opts *handlers.GenerateOptions) error {
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
func promptForResources(opts *handlers.GenerateOptions) error {
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
func promptForCapabilities(opts *handlers.GenerateOptions) error {
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
