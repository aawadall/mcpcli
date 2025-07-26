package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/aawadall/mcpcli/internal/core"
	"github.com/fatih/color"
)

// PythonGenerator implements the Generator interface for Python projects
type PythonGenerator struct{}

func NewPythonGenerator() *PythonGenerator { return &PythonGenerator{} }

func (g *PythonGenerator) GetLanguage() string { return "python" }

func (g *PythonGenerator) GetSupportedTransports() []string { return []string{"stdio"} }

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

func (g *PythonGenerator) generateFromTemplates(output string, data *core.TemplateData) error {
	templates := map[string]string{
		"templates/python/stdio/src/main.py.tmpl":               "src/main.py",
		"templates/python/stdio/src/handlers/mcp.py.tmpl":       "src/handlers/mcp.py",
		"templates/python/stdio/src/resources/registry.py.tmpl": "src/resources/registry.py",
		"templates/python/stdio/README.md.tmpl":                 "README.md",
		"templates/python/stdio/configs/mcp-config.json.tmpl":   "configs/mcp-config.json",
		"templates/python/stdio/examples/example.py.tmpl":       "examples/example.py",
	}
	for tPath, outPath := range templates {
		if err := g.generateTemplate(tPath, filepath.Join(output, outPath), data); err != nil {
			return err
		}
	}
	if data.Config.Docker {
		dockerTemps := map[string]string{
			"templates/python/stdio/Dockerfile.tmpl":   "Dockerfile",
			"templates/python/stdio/dockerignore.tmpl": ".dockerignore",
		}
		for tPath, outPath := range dockerTemps {
			if err := g.generateTemplate(tPath, filepath.Join(output, outPath), data); err != nil {
				return err
			}
		}
	}

	for _, tool := range data.Config.Tools {
		td := struct{ Tool core.Tool }{Tool: tool}
		file := filepath.Join(output, "src/tools", tool.Name+".py")
		if err := g.generateTemplate("templates/python/stdio/src/tools/tool.py.tmpl", file, td); err != nil {
			return err
		}
	}

	for _, res := range data.Config.Resources {
		rd := struct{ Resource core.Resource }{Resource: res}
		file := filepath.Join(output, "src/resources", res.Name+".py")
		if err := g.generateTemplate("templates/python/stdio/src/resources/resource.py.tmpl", file, rd); err != nil {
			return err
		}
	}

	for _, cap := range data.Config.Capabilities {
		cd := struct{ Capability core.Capability }{Capability: cap}
		file := filepath.Join(output, "src/capabilities", cap.Name+".py")
		if err := g.generateTemplate("templates/python/stdio/src/capabilities/capability.py.tmpl", file, cd); err != nil {
			return err
		}
	}

	return nil
}

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
