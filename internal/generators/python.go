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

// PythonGenerator implements the Generator interface for Python projects.
type PythonGenerator struct{}

// NewPythonGenerator returns a new instance of PythonGenerator.
func NewPythonGenerator() *PythonGenerator { return &PythonGenerator{} }

// GetLanguage returns the generator language.
func (g *PythonGenerator) GetLanguage() string { return "python" }

// GetSupportedTransports lists the transports supported by the generator.
func (g *PythonGenerator) GetSupportedTransports() []string { return []string{"stdio"} }

// Generate scaffolds a Python project using the provided configuration.
func (g *PythonGenerator) Generate(config *core.ProjectConfig) error {
	data := config.GetTemplateData()
	if err := g.createDirectoryStructure(config.Output); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}
	if err := g.generateFromTemplates(config.Output, data); err != nil {
		return fmt.Errorf("failed to generate from templates: %w", err)
	}
	return nil
}

// createDirectoryStructure initializes the folder layout for the project.
func (g *PythonGenerator) createDirectoryStructure(output string) error {
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

// generateFromTemplates renders all static and dynamic templates for the project.
func (g *PythonGenerator) generateFromTemplates(output string, data *core.TemplateData) error {
	templates, err := tmp.BaseTemplateMap(g.GetLanguage(), data)
	if err != nil {
		return err
	}
	if err := g.generateTemplateMap(output, templates, data); err != nil {
		return err
	}

	// Docker templates are included in BaseTemplateMap

	toolConv := func(t core.Tool) (string, interface{}) {
		return t.Name, struct{ Tool core.Tool }{Tool: t}
	}
	if err := generateEntities(g, output, "src/tools", tmp.ToolTemplate(g.GetLanguage()), data.Config.Tools, toolConv); err != nil {
		return err
	}

	resConv := func(r core.Resource) (string, interface{}) {
		return r.Name, struct{ Resource core.Resource }{Resource: r}
	}
	if err := generateEntities(g, output, "src/resources", tmp.ResourceTemplate(g.GetLanguage()), data.Config.Resources, resConv); err != nil {
		return err
	}

	capConv := func(c core.Capability) (string, interface{}) {
		return c.Name, struct{ Capability core.Capability }{Capability: c}
	}
	if err := generateEntities(g, output, "src/capabilities", tmp.CapabilityTemplate(g.GetLanguage()), data.Config.Capabilities, capConv); err != nil {
		return err
	}

	return nil
}

// generateTemplate renders a single template file to the destination path.
func (g *PythonGenerator) generateTemplate(tPath, outPath string, data interface{}) error {
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

// generateTemplateMap iterates over a map of templates and generates each file.
func (g *PythonGenerator) generateTemplateMap(output string, templates map[string]string, data interface{}) error {
	for tPath, outPath := range templates {
		if err := g.generateTemplate(tPath, filepath.Join(output, outPath), data); err != nil {
			return err
		}
	}
	return nil
}

// generateEntities iterates over the given items, converts each to template data
// using conv, and writes the resulting file to the specified directory.
// The file is named after the item with a .py extension.
func generateEntities[T any](g *PythonGenerator, output, dir, tmpl string, items []T, conv func(T) (string, interface{})) error {
	for _, item := range items {
		name, d := conv(item)
		file := filepath.Join(output, dir, name+".py")
		if err := g.generateTemplate(tmpl, file, d); err != nil {
			return err
		}
	}
	return nil
}
