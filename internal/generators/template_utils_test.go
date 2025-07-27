package generators

import (
	"testing"

	"github.com/aawadall/mcpcli/internal/core"
	tmp "github.com/aawadall/mcpcli/internal/generators/templates"
)

func TestBaseTemplateMap_Go(t *testing.T) {
	cfg := &core.ProjectConfig{}
	data := cfg.GetTemplateData()
	m, err := tmp.BaseTemplateMap("go", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["templates/go/stdio/go.mod.tmpl"] == "" {
		t.Errorf("expected go.mod template mapping")
	}
}

func TestBaseTemplateMap_Docker(t *testing.T) {
	cfg := &core.ProjectConfig{Docker: true}
	data := cfg.GetTemplateData()
	m, err := tmp.BaseTemplateMap("python", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m["templates/python/stdio/Dockerfile.tmpl"]; !ok {
		t.Errorf("expected docker template when docker enabled")
	}
}

func TestBaseTemplateMap_Invalid(t *testing.T) {
	cfg := &core.ProjectConfig{}
	data := cfg.GetTemplateData()
	if _, err := tmp.BaseTemplateMap("invalid", data); err == nil {
		t.Error("expected error for invalid language")
	}
}

func TestTemplateHelpers(t *testing.T) {
	if tmp.ToolTemplate("go") == "" || tmp.ResourceTemplate("go") == "" || tmp.CapabilityTemplate("go") == "" {
		t.Error("expected template helper paths for go")
	}
}
