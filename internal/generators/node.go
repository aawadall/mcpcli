package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/aawadall/mcpcli/internal/core"
	"github.com/fatih/color"
)

// NodeGenerator implements the Generator interface for Node.js projects
type NodeGenerator struct{}

func NewNodeGenerator() *NodeGenerator {
	return &NodeGenerator{}
}

func (g *NodeGenerator) GetLanguage() string {
	return "javascript"
}

func (g *NodeGenerator) GetSupportedTransports() []string {
	return []string{"stdio"}
}

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

func (g *NodeGenerator) generateFromTemplates(output string, data *core.TemplateData) error {
	templates := map[string]string{
		"templates/node/stdio/package.json.tmpl":              "package.json",
		"templates/node/stdio/src/index.js.tmpl":              "src/index.js",
		"templates/node/stdio/src/handlers/mcp.js.tmpl":       "src/handlers/mcp.js",
		"templates/node/stdio/src/resources/registry.js.tmpl": "src/resources/registry.js",
		"templates/node/stdio/README.md.tmpl":                 "README.md",
		"templates/node/stdio/configs/mcp-config.json.tmpl":   "configs/mcp-config.json",
		"templates/node/stdio/examples/example.js.tmpl":       "examples/example.js",
	}

	for tPath, outPath := range templates {
		if err := g.generateTemplate(tPath, filepath.Join(output, outPath), data); err != nil {
			return err
		}
	}

	if data.Config.Docker {
		dockerTemps := map[string]string{
			"templates/node/stdio/Dockerfile.tmpl":   "Dockerfile",
			"templates/node/stdio/dockerignore.tmpl": ".dockerignore",
		}
		for tPath, outPath := range dockerTemps {
			if err := g.generateTemplate(tPath, filepath.Join(output, outPath), data); err != nil {
				return err
			}
		}
	}

	for _, tool := range data.Config.Tools {
		td := struct {
			Tool core.Tool
		}{Tool: tool}
		file := filepath.Join(output, "src/tools", tool.Name+".js")
		if err := g.generateTemplate("templates/node/stdio/src/tools/tool.js.tmpl", file, td); err != nil {
			return err
		}
	}

	for _, res := range data.Config.Resources {
		rd := struct {
			Resource core.Resource
		}{Resource: res}
		file := filepath.Join(output, "src/resources", res.Name+".js")
		if err := g.generateTemplate("templates/node/stdio/src/resources/resource.js.tmpl", file, rd); err != nil {
			return err
		}
	}

	for _, cap := range data.Config.Capabilities {
		cd := struct {
			Capability core.Capability
		}{Capability: cap}
		file := filepath.Join(output, "src/capabilities", cap.Name+".js")
		if err := g.generateTemplate("templates/node/stdio/src/capabilities/capability.js.tmpl", file, cd); err != nil {
			return err
		}
	}

	return nil
}

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
