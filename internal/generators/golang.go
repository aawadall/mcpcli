package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/aawadall/mcpcli/internal/core"
	"github.com/aawadall/mcpcli/internal/generators/templates"
	"github.com/fatih/color"
)

// GoGenerator implements the Generator interface for Go projects
type GoGenerator struct{}

func NewGolangGenerator() *GoGenerator {
	return &GoGenerator{}
}

// Generate creates a Go project structure based on the provided configuration
func (g *GoGenerator) Generate(config *core.ProjectConfig) error {
	templateData := config.GetTemplateData()

	// Create directory structure
	if err := g.createDirectoryStructure(config.Output); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Generate files using templates
	if err := g.generateFromTemplates(config.Output, templateData); err != nil {
		return fmt.Errorf("failed to generate from templates: %w", err)
	}

	return nil
}

// GetLanguage returns the language name
func (g *GoGenerator) GetLanguage() string {
	return "go"
}

// GetSupportedTransports returns the list of supported transports for Go
func (g *GoGenerator) GetSupportedTransports() []string {
	return []string{"stdio", "rest", "websocket"}
}

// createDirectoryStructure creates the project directory structure
func (g *GoGenerator) createDirectoryStructure(output string) error {
	dirs := []string{
		"cmd/server",
		"internal/handlers",
		"internal/resources",
		"internal/tools",
		"internal/capabilities",
		"pkg/mcp",
		"examples",
		"configs",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(output, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

// generateFromTemplates generates all files using the template system
func (g *GoGenerator) generateFromTemplates(output string, data *core.TemplateData) error {
	// Define template files and their output paths
	templates, err := templates.BaseTemplateMap(g.GetLanguage(), data)
	if err != nil {
		return err
	}

	for templatePath, outputPath := range templates {
		if err := g.generateTemplate(templatePath, filepath.Join(output, outputPath), data); err != nil {
			return fmt.Errorf("failed to generate %s: %w", outputPath, err)
		}
	}

	// Generate a dedicated Go file for each tool
	for _, tool := range data.Config.Tools {
		toolData := struct {
			ModuleName string
			Tool       core.Tool
		}{
			ModuleName: data.ModuleName,
			Tool:       tool,
		}
		fileName := filepath.Join(output, "internal/tools", tool.Name+".go")
		tmplPath := templates.ToolTemplate(g.GetLanguage())
		if err := g.generateTemplate(tmplPath, fileName, toolData); err != nil {
			return fmt.Errorf("failed to generate tool file for %s: %w", tool.Name, err)
		}
	}

	// Generate a dedicated Go file for each resource
	for _, resource := range data.Config.Resources {
		resourceData := struct {
			ModuleName string
			Resource   core.Resource
		}{
			ModuleName: data.ModuleName,
			Resource:   resource,
		}
		fileName := filepath.Join(output, "internal/resources", resource.Name+".go")
		tmplPath := templates.ResourceTemplate(g.GetLanguage())
		if err := g.generateTemplate(tmplPath, fileName, resourceData); err != nil {
			return fmt.Errorf("failed to generate resource file for %s: %w", resource.Name, err)
		}
	}

	// Generate a dedicated Go file for each capability
	for _, capability := range data.Config.Capabilities {
		capabilityData := struct {
			ModuleName string
			Capability core.Capability
		}{
			ModuleName: data.ModuleName,
			Capability: capability,
		}
		fileName := filepath.Join(output, "internal/capabilities", capability.Name+".go")
		tmplPath := templates.CapabilityTemplate(g.GetLanguage())
		if err := g.generateTemplate(tmplPath, fileName, capabilityData); err != nil {
			return fmt.Errorf("failed to generate capability file for %s: %w", capability.Name, err)
		}
	}

	return nil
}

// generateTemplate generates a single file from a template
func (g *GoGenerator) generateTemplate(templatePath, outputPath string, data interface{}) error {
	// Read template content
	templateContent, err := TemplatesFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	color.Green("✅ Created file: %s", outputPath)
	return nil
}
