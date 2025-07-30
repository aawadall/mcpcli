package handlers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/aawadall/mcpcli/internal/core"
)

func writeConfig(t *testing.T, dir string) string {
	cfg := &core.MCPConfig{Name: "test", Version: "0.4.2", Transport: core.Transport{Type: "stdio", Options: map[string]any{"command": "true"}}}
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "cfg.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadMCPConfig_Invalid(t *testing.T) {
	tmp := t.TempDir()
	bad := filepath.Join(tmp, "bad.json")
	os.WriteFile(bad, []byte("{"), 0644)
	if _, err := LoadMCPConfig(bad); err == nil {
		t.Fatal("expected error for invalid json")
	}
}

func TestLoadMCPConfig_Valid(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := writeConfig(t, tmp)
	cfg, err := LoadMCPConfig(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "test" {
		t.Fatalf("expected name 'test', got %s", cfg.Name)
	}
}

func TestRunTests_BadCommand(t *testing.T) {
	opts := &TestOptions{TestAll: true}
	cfg := &core.MCPConfig{Name: "bad", Transport: core.Transport{Type: "stdio", Options: map[string]any{"command": "nonexistent-cmd"}}}
	if err := RunTests(opts, cfg); err == nil {
		t.Fatal("expected error for bad command")
	}
}

func TestRunTests_NoServer(t *testing.T) {
	opts := &TestOptions{TestAll: true}
	cfg := &core.MCPConfig{Name: "none", Transport: core.Transport{Type: "stdio", Options: map[string]any{"command": "true"}}}
	if err := RunTests(opts, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestRunTests_WithScript(t *testing.T) {
	opts := &TestOptions{ScriptFile: "script.txt"}
	cfg := &core.MCPConfig{Name: "s", Transport: core.Transport{Type: "stdio", Options: map[string]any{"command": "true"}}}
	if err := RunTests(opts, cfg); err != nil {
		t.Fatalf("RunTests with script failed: %v", err)
	}
}

func TestLoadMCPConfig_Project(t *testing.T) {
	pc := core.NewProjectConfig()
	pc.Name = "proj"
	data, _ := json.Marshal(pc)
	dir := t.TempDir()
	path := filepath.Join(dir, "proj.json")
	os.WriteFile(path, data, 0644)
	cfg, err := LoadMCPConfig(path)
	if err != nil {
		t.Fatalf("load project: %v", err)
	}
	if cfg.Name != "proj" {
		t.Fatalf("expected proj name, got %s", cfg.Name)
	}
}
