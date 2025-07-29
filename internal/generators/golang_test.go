package generators

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aawadall/mcpcli/internal/core"
)

func TestGoGenerator_GetLanguage(t *testing.T) {
	g := NewGolangGenerator()
	if g.GetLanguage() != "go" {
		t.Errorf("expected language 'go', got '%s'", g.GetLanguage())
	}
}

func TestGoGenerator_GetSupportedTransports(t *testing.T) {
	g := NewGolangGenerator()
	transports := g.GetSupportedTransports()
	if len(transports) == 0 {
		t.Error("expected at least one supported transport")
	}
}

func TestGoGenerator_createDirectoryStructure(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGolangGenerator()
	err := g.createDirectoryStructure(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dirs := []string{
		"cmd/server",
		"internal/handlers",
		"internal/resources",
		"internal/tools",
		"pkg/mcp",
		"examples",
		"configs",
	}
	for _, dir := range dirs {
		path := filepath.Join(tmpDir, dir)
		if stat, err := os.Stat(path); err != nil || !stat.IsDir() {
			t.Errorf("expected directory %s to exist", path)
		}
	}
}

func TestGoGenerator_generateTemplate_invalidPath(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGolangGenerator()
	data := &core.TemplateData{}
	err := g.generateTemplate("nonexistent.tmpl", filepath.Join(tmpDir, "out.go"), data)
	if err == nil {
		t.Error("expected error for nonexistent template path")
	}
}

func TestGoGenerator_GenerateWithExtras(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &core.ProjectConfig{
		Name:         "extras",
		Language:     "go",
		Transport:    "stdio",
		Output:       tmpDir,
		Tools:        []core.Tool{{Name: "Hammer"}},
		Resources:    []core.Resource{{Name: "Nail", Type: string(core.ResourceTypeFilesystem)}},
		Capabilities: []core.Capability{{Name: "Build"}},
	}

	g := NewGolangGenerator()
	if err := g.Generate(cfg); err != nil {
		t.Fatalf("unexpected error generating project: %v", err)
	}

	expected := []string{
		filepath.Join(tmpDir, "internal", "tools", "Hammer.go"),
		filepath.Join(tmpDir, "internal", "resources", "Nail.go"),
		filepath.Join(tmpDir, "internal", "capabilities", "Build.go"),
	}
	for _, f := range expected {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected file %s to exist, got %v", f, err)
		}
	}
}
