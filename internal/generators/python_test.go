package generators

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aawadall/mcpcli/internal/core"
)

func TestPythonGenerator_GetLanguage(t *testing.T) {
	g := NewPythonGenerator()
	if g.GetLanguage() != "python" {
		t.Errorf("expected language 'python', got '%s'", g.GetLanguage())
	}
}

func TestPythonGenerator_GetSupportedTransports(t *testing.T) {
	g := NewPythonGenerator()
	tr := g.GetSupportedTransports()
	if len(tr) == 0 {
		t.Error("expected at least one supported transport")
	}
	if tr[0] != "stdio" {
		t.Errorf("expected transport 'stdio', got '%s'", tr[0])
	}
}

func TestPythonGenerator_GenerateBasic(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &core.ProjectConfig{
		Name:      "testpython",
		Language:  "python",
		Transport: "stdio",
		Output:    tmpDir,
	}

	g := NewPythonGenerator()
	if err := g.Generate(cfg); err != nil {
		t.Fatalf("unexpected error generating project: %v", err)
	}

	expected := []string{
		filepath.Join(tmpDir, "src", "main.py"),
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
