package generators

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aawadall/mcpcli/internal/core"
)

func TestNodeGenerator_GetLanguage(t *testing.T) {
	g := NewNodeGenerator()
	if g.GetLanguage() != "javascript" {
		t.Errorf("expected language 'javascript', got '%s'", g.GetLanguage())
	}
}

func TestNodeGenerator_GetSupportedTransports(t *testing.T) {
	g := NewNodeGenerator()
	transports := g.GetSupportedTransports()
	if len(transports) == 0 {
		t.Error("expected at least one supported transport")
	}
	if transports[0] != "stdio" {
		t.Errorf("expected transport 'stdio', got '%s'", transports[0])
	}
}

func TestNodeGenerator_GenerateBasic(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &core.ProjectConfig{
		Name:      "testnode",
		Language:  "javascript",
		Transport: "stdio",
		Output:    tmpDir,
	}

	g := NewNodeGenerator()
	if err := g.Generate(cfg); err != nil {
		t.Fatalf("unexpected error generating project: %v", err)
	}

	expected := []string{
		filepath.Join(tmpDir, "package.json"),
		filepath.Join(tmpDir, "src", "index.js"),
	}

	for _, f := range expected {
		info, err := os.Stat(f)
		if err != nil {
			t.Errorf("expected file %s to exist, but got error: %v", f, err)
			continue
		}
		if info.IsDir() {
			t.Errorf("expected file %s to be a file, but it is a directory", f)
		}
	}
}
