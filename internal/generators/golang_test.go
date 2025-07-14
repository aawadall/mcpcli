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

// Note: More comprehensive tests for generateFromTemplates and Generate would require
// integration-style tests with actual template files and a mock config.
