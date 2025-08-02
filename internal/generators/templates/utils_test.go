package generators

import (
	"testing"

	"github.com/aawadall/mcpcli/internal/core"
)

func TestBaseTemplateMapUnsupported(t *testing.T) {
	_, err := BaseTemplateMap("bad", &core.TemplateData{})
	if err == nil {
		t.Fatal("expected error for unsupported language")
	}
}

func TestTemplateHelpers(t *testing.T) {
	data := &core.TemplateData{Config: &core.ProjectConfig{Docker: true}}
	langs := []string{"go", "javascript", "python", "java"}
	for _, lang := range langs {
		m, err := BaseTemplateMap(lang, data)
		if err != nil {
			t.Fatalf("map for %s: %v", lang, err)
		}
		if len(m) == 0 {
			t.Fatalf("expected templates for %s", lang)
		}
		if ToolTemplate(lang) == "" {
			t.Fatalf("tool template empty for %s", lang)
		}
		if ResourceTemplate(lang) == "" {
			t.Fatalf("resource template empty for %s", lang)
		}
		if CapabilityTemplate(lang) == "" {
			t.Fatalf("capability template empty for %s", lang)
		}
	}
}

func TestBaseTemplateMapRestTransport(t *testing.T) {
	cfg := &core.ProjectConfig{Transport: "rest"}
	data := cfg.GetTemplateData()
	expected := map[string]string{
		"go":         "templates/go/http/cmd/server/main.go.tmpl",
		"javascript": "templates/node/http/src/index.js.tmpl",
		"python":     "templates/python/http/src/main.py.tmpl",
		"java":       "templates/java/http/src/main/java/Main.java.tmpl",
	}
	for lang, path := range expected {
		m, err := BaseTemplateMap(lang, data)
		if err != nil {
			t.Fatalf("map for %s: %v", lang, err)
		}
		if _, ok := m[path]; !ok {
			t.Errorf("expected %s to map to %s", lang, path)
		}
	}
}
