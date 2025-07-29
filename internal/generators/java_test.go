package generators

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aawadall/mcpcli/internal/core"
)

func TestJavaGenerator_GetLanguage(t *testing.T) {
	g := NewJavaGenerator()
	if g.GetLanguage() != "java" {
		t.Errorf("expected language 'java', got '%s'", g.GetLanguage())
	}
}

func TestJavaGenerator_GetSupportedTransports(t *testing.T) {
	g := NewJavaGenerator()
	transports := g.GetSupportedTransports()
	if len(transports) == 0 {
		t.Error("expected at least one supported transport")
	}
	if transports[0] != "stdio" {
		t.Errorf("expected transport 'stdio', got '%s'", transports[0])
	}
}

func TestJavaGenerator_GenerateBasic(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &core.ProjectConfig{
		Name:      "testjava",
		Language:  "java",
		Transport: "stdio",
		Output:    tmpDir,
	}

	g := NewJavaGenerator()
	if err := g.Generate(cfg); err != nil {
		t.Fatalf("unexpected error generating project: %v", err)
	}

	pkg := cfg.GetTemplateData().PackageName
	expected := []string{
		filepath.Join(tmpDir, "pom.xml"),
		filepath.Join(tmpDir, "src", "main", "java", pkg, "Main.java"),
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

func TestJavaGenerator_GenerateWithExtras(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &core.ProjectConfig{
		Name:         "extrajava",
		Language:     "java",
		Transport:    "stdio",
		Output:       tmpDir,
		Tools:        []core.Tool{{Name: "Hammer"}},
		Resources:    []core.Resource{{Name: "Nail"}},
		Capabilities: []core.Capability{{Name: "Build"}},
	}

	g := NewJavaGenerator()
	if err := g.Generate(cfg); err != nil {
		t.Fatalf("unexpected error generating project: %v", err)
	}

	pkg := cfg.GetTemplateData().PackageName
	expected := []string{
		filepath.Join(tmpDir, "src", "main", "java", pkg, "tools", "Hammer.java"),
		filepath.Join(tmpDir, "src", "main", "java", pkg, "resources", "Nail.java"),
		filepath.Join(tmpDir, "src", "main", "java", pkg, "capabilities", "Build.java"),
	}
	for _, f := range expected {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected file %s to exist, got %v", f, err)
		}
	}
}
