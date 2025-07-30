package generators

import (
	"os"
	"path/filepath"
	"strings"
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

func TestPythonGenerator_GenerateWithExtras(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &core.ProjectConfig{
		Name:         "extras",
		Language:     "python",
		Transport:    "stdio",
		Output:       tmpDir,
		Tools:        []core.Tool{{Name: "tool"}},
		Resources:    []core.Resource{{Name: "res"}},
		Capabilities: []core.Capability{{Name: "cap"}},
	}

	g := NewPythonGenerator()
	if err := g.Generate(cfg); err != nil {
		t.Fatalf("unexpected error generating project: %v", err)
	}

	expected := []string{
		filepath.Join(tmpDir, "src", "tools", "tool.py"),
		filepath.Join(tmpDir, "src", "resources", "res.py"),
		filepath.Join(tmpDir, "src", "capabilities", "cap.py"),
	}

	for _, f := range expected {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected file %s to exist, got %v", f, err)
		}
	}
}

func TestPythonCreateDirectoryStructure_Error(t *testing.T) {
	tmpDir := t.TempDir()
	// create a file where a directory is expected
	if err := os.WriteFile(filepath.Join(tmpDir, "src"), []byte("file"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	g := NewPythonGenerator()
	if err := g.createDirectoryStructure(tmpDir); err == nil {
		t.Fatal("expected error when creating directories over existing file")
	}
}

func TestPythonGenerateTemplate_ReadError(t *testing.T) {
	g := NewPythonGenerator()
	tmpDir := t.TempDir()
	err := g.generateTemplate("missing.tmpl", filepath.Join(tmpDir, "out.py"), nil)
	if err == nil || !strings.Contains(err.Error(), "failed to read template") {
		t.Fatalf("expected template read error, got %v", err)
	}
}

func TestPythonGenerateTemplate_CreateError(t *testing.T) {
	g := NewPythonGenerator()
	tmpDir := t.TempDir()
	tPath := "templates/python/stdio/src/main.py.tmpl"
	outPath := filepath.Join(tmpDir, "nope", "file.py")
	err := g.generateTemplate(tPath, outPath, nil)
	if err == nil || !strings.Contains(err.Error(), "failed to create file") {
		t.Fatalf("expected file creation error, got %v", err)
	}
}

func TestPythonGenerateEntities_Error(t *testing.T) {
	g := NewPythonGenerator()
	tmpDir := t.TempDir()
	items := []core.Tool{{Name: "bad"}}
	conv := func(t core.Tool) (string, interface{}) { return t.Name, struct{ Tool core.Tool }{t} }
	err := generateEntities(g, tmpDir, "src/tools", "missing.tmpl", items, conv)
	if err == nil || !strings.Contains(err.Error(), "failed to read template") {
		t.Fatalf("expected template read error, got %v", err)
	}
}

func TestPythonGenerateTemplateMap_Error(t *testing.T) {
	g := NewPythonGenerator()
	tmpDir := t.TempDir()
	templates := map[string]string{"missing.tmpl": "out.py"}
	err := g.generateTemplateMap(tmpDir, templates, nil)
	if err == nil || !strings.Contains(err.Error(), "failed to read template") {
		t.Fatalf("expected template map error, got %v", err)
	}
}

func TestPythonGenerateFromTemplates_Error(t *testing.T) {
	g := NewPythonGenerator()
	tmpDir := t.TempDir()
	// create a file "src" to block directory creation
	if err := os.WriteFile(filepath.Join(tmpDir, "src"), []byte("file"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	cfg := &core.ProjectConfig{Name: "bad", Language: "python", Transport: "stdio", Output: tmpDir}
	data := cfg.GetTemplateData()
	err := g.generateFromTemplates(tmpDir, data)
	if err == nil || !strings.Contains(err.Error(), "failed to create file") {
		t.Fatalf("expected file creation error, got %v", err)
	}
}

func TestPythonGenerator_Generate_DirectoryError(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "src"), []byte("file"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	cfg := &core.ProjectConfig{Name: "badpython", Language: "python", Transport: "stdio", Output: tmpDir}
	g := NewPythonGenerator()
	err := g.Generate(cfg)
	if err == nil || !strings.Contains(err.Error(), "failed to create directory structure") {
		t.Fatalf("expected directory structure error, got %v", err)
	}
}
