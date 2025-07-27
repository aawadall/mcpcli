package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/aawadall/mcpcli/internal/core"
	tmp "github.com/aawadall/mcpcli/internal/generators/templates"
	"github.com/fatih/color"
)

// JavaGenerator implements the Generator interface for Java projects
type JavaGenerator struct{}

func NewJavaGenerator() *JavaGenerator { return &JavaGenerator{} }

func (g *JavaGenerator) GetLanguage() string { return "java" }

func (g *JavaGenerator) GetSupportedTransports() []string { return []string{"stdio"} }

func (g *JavaGenerator) Generate(config *core.ProjectConfig) error {
	data := config.GetTemplateData()
	if err := g.createDirectoryStructure(config.Output, data.PackageName); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}
	if err := g.generateFromTemplates(config.Output, data); err != nil {
		return fmt.Errorf("failed to generate from templates: %w", err)
	}
	return nil
}

func (g *JavaGenerator) createDirectoryStructure(output, pkg string) error {
	pkgPath := filepath.Join(strings.Split(pkg, ".")...)
	dirs := []string{
		filepath.Join("src", "main", "java", pkgPath, "handlers"),
		filepath.Join("src", "main", "java", pkgPath, "resources"),
		filepath.Join("src", "main", "java", pkgPath, "tools"),
		filepath.Join("src", "main", "java", pkgPath, "capabilities"),
		filepath.Join("src", "main", "java", pkgPath),
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

func (g *JavaGenerator) generateFromTemplates(output string, data *core.TemplateData) error {
	pkgPath := filepath.Join(strings.Split(data.PackageName, ".")...)
	templates, err := tmp.BaseTemplateMap(g.GetLanguage(), data)
	if err != nil {
		return err
	}
	for tPath, outPath := range templates {
		if err := g.generateTemplate(tPath, filepath.Join(output, outPath), data); err != nil {
			return err
		}
	}
	for _, tool := range data.Config.Tools {
		td := struct {
			PackageName string
			Tool        core.Tool
		}{PackageName: data.PackageName, Tool: tool}
		file := filepath.Join(output, "src", "main", "java", pkgPath, "tools", tool.Name+".java")
		if err := g.generateTemplate(tmp.ToolTemplate(g.GetLanguage()), file, td); err != nil {
			return err
		}
	}
	for _, res := range data.Config.Resources {
		rd := struct {
			PackageName string
			Resource    core.Resource
		}{PackageName: data.PackageName, Resource: res}
		file := filepath.Join(output, "src", "main", "java", pkgPath, "resources", res.Name+".java")
		if err := g.generateTemplate(tmp.ResourceTemplate(g.GetLanguage()), file, rd); err != nil {
			return err
		}
	}
	for _, cap := range data.Config.Capabilities {
		cd := struct {
			PackageName string
			Capability  core.Capability
		}{PackageName: data.PackageName, Capability: cap}
		file := filepath.Join(output, "src", "main", "java", pkgPath, "capabilities", cap.Name+".java")
		if err := g.generateTemplate(tmp.CapabilityTemplate(g.GetLanguage()), file, cd); err != nil {
			return err
		}
	}
	return nil
}

func (g *JavaGenerator) generateTemplate(tPath, outPath string, data interface{}) error {
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
