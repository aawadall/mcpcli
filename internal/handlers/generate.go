package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aawadall/mcpcli/internal/core"
	"github.com/aawadall/mcpcli/internal/generators"
)

// GenerateOptions holds all configurable parameters for project generation.
type GenerateOptions struct {
	Name      string
	Language  string
	Transport string
	Docker    bool
	Examples  bool
	Output    string
	// Interactive indicates if prompts should be shown. It is ignored by the generator.
	Interactive  bool
	Force        bool
	Tools        []core.Tool
	Resources    []core.Resource
	Capabilities []core.Capability
}

// ValidateGenerateOptions checks if the provided options are valid.
func ValidateGenerateOptions(opts *GenerateOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("project name is required")
	}
	validLanguages := []string{"golang", "javascript", "java", "python"}
	if !contains(validLanguages, opts.Language) {
		return fmt.Errorf("invalid language: %s, valid options are: %v", opts.Language, validLanguages)
	}
	validTransports := []string{"stdio", "rest", "websocket"}
	if !contains(validTransports, opts.Transport) {
		return fmt.Errorf("invalid transport: %s, valid options are: %v", opts.Transport, validTransports)
	}
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

// GenerateProject creates the project structure based on the provided options.
func GenerateProject(opts *GenerateOptions) error {
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
