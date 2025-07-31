package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/aawadall/mcpcli/internal/core"
	tmp "github.com/aawadall/mcpcli/internal/generators/templates"
	"github.com/fatih/color"
)

// NodeGenerator implements the Generator interface for Node.js projects.
type NodeGenerator struct{}

// NewNodeGenerator creates a new NodeGenerator instance.
func NewNodeGenerator() *NodeGenerator {
	return &NodeGenerator{}
}

// GetLanguage returns the language identifier for the generator.
func (g *NodeGenerator) GetLanguage() string {
	return "javascript"
}

// GetSupportedTransports lists the supported transport mechanisms.
func (g *NodeGenerator) GetSupportedTransports() []string {
	return []string{"stdio", "rest", "websocket"}
}

// Generate scaffolds a Node.js project using the provided configuration.
func (g *NodeGenerator) Generate(config *core.ProjectConfig) error {
	data := config.GetTemplateData()
	if err := g.createDirectoryStructure(config.Output); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}
	if err := g.generateFromTemplates(config.Output, data); err != nil {
		return fmt.Errorf("failed to generate from templates: %w", err)
	}
	return nil
}

// createDirectoryStructure creates the required directories for the project.
func (g *NodeGenerator) createDirectoryStructure(output string) error {
	dirs := []string{
		"src/handlers",
		"src/resources",
		"src/tools",
		"src/capabilities",
		"examples",
		"configs",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(output, d), 0755); err != nil {
			return err
		}
	}
	return nil
}

// generateFromTemplates writes all base templates and additional items.
func (g *NodeGenerator) generateFromTemplates(output string, data *core.TemplateData) error {
	templates, err := tmp.BaseTemplateMap(g.GetLanguage(), data)
	if err != nil {
		return err
	}

	// Generate the main project files
	for tPath, outPath := range templates {
		if err := g.generateTemplate(tPath, filepath.Join(output, outPath), data); err != nil {
			return err
		}
	}

	// Generate tools
	if err := generateItems(g, output, "src/tools", data.Config.Tools, tmp.ToolTemplate(g.GetLanguage())); err != nil {
		return err
	}

	// generate resources
	if err := generateItems(g, output, "src/resources", data.Config.Resources, tmp.ResourceTemplate(g.GetLanguage())); err != nil {
		return err
	}

	// generate capabilities
	if err := generateItems(g, output, "src/capabilities", data.Config.Capabilities, tmp.CapabilityTemplate(g.GetLanguage())); err != nil {
		return err
	}

	return nil
}

// generateItems creates files for a slice of items that implement the GetName() method.
// It uses the specified template and writes the generated files to the given subdirectory
// within the output directory. Each item's name, as returned by GetName(), is used to
// determine the filename.
// generateItems iterates over items and renders a template for each one.
func generateItems[T interface{ GetName() string }](g *NodeGenerator, output, subDir string, items []T, templatePath string) error {
	for _, item := range items {
		// create wrapper struct with the item
		data := struct {
			Item T
		}{Item: item}

		// generate the file path
		filePath := filepath.Join(output, subDir, item.GetName()+".js")
		if err := g.generateTemplate(templatePath, filePath, data); err != nil {
			return err
		}
	}
	return nil
}

// generateTemplate renders a single template file with the provided data
// and writes the output to the specified path.
func (g *NodeGenerator) generateTemplate(tPath, outPath string, data interface{}) error {
	content, err := TemplatesFS.ReadFile(tPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tPath, err)
	}
	tmpl, err := template.New(filepath.Base(tPath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", tPath, err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outPath, err)
	}
	defer f.Close()
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", tPath, err)
	}
	color.Green("✅ Created file: %s", outPath)
	return nil
}
