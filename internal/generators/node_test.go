package generators

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aawadall/mcpcli/internal/core"
	"github.com/stretchr/testify/assert"
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

// Mock types for testing
type mockTool struct {
	Name string
}

func (m mockTool) GetName() string {
	return m.Name
}

type mockResource struct {
	Name string
}

func (m mockResource) GetName() string {
	return m.Name
}

type mockCapability struct {
	Name string
}

func (m mockCapability) GetName() string {
	return m.Name
}

func TestGenerateItems_Logic_Only(t *testing.T) {
	// Test the path generation logic
	tools := []mockTool{{Name: "testTool"}}
	output := "/project"
	subDir := "src/tools"

	for _, tool := range tools {
		expectedPath := filepath.Join(output, subDir, tool.GetName()+".js")
		assert.Equal(t, "/project/src/tools/testTool.js", expectedPath)
	}

	// Test the data structure creation
	tool := mockTool{Name: "myTool"}
	data := struct {
		Item mockTool
	}{Item: tool}

	assert.Equal(t, "myTool", data.Item.GetName())
}

func TestTemplateDataStructure(t *testing.T) {
	// Test that the wrapper struct works correctly
	tool := mockTool{Name: "hammer"}
	resource := mockResource{Name: "database"}
	capability := mockCapability{Name: "read"}

	toolData := struct{ Item mockTool }{Item: tool}
	resourceData := struct{ Item mockResource }{Item: resource}
	capabilityData := struct{ Item mockCapability }{Item: capability}

	assert.Equal(t, "hammer", toolData.Item.GetName())
	assert.Equal(t, "database", resourceData.Item.GetName())
	assert.Equal(t, "read", capabilityData.Item.GetName())
}
