package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aawadall/mcpcli/internal/core"
)

func writeTempConfig(t *testing.T, dir string) string {
	cfg := &core.MCPConfig{
		Name:      "test",
		Version:   "0.4.1",
		Transport: core.Transport{Type: "stdio"},
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestTestCmd_FlagParsingAndExecution(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := writeTempConfig(t, tmpDir)
	script := filepath.Join(tmpDir, "script.txt")

	cmd := NewTestCmd()
	cmd.SetArgs([]string{"--config", cfg, "--script", script})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	confVal, _ := cmd.Flags().GetString("config")
	if confVal != cfg {
		t.Errorf("expected config %s, got %s", cfg, confVal)
	}
	scriptVal, _ := cmd.Flags().GetString("script")
	if scriptVal != script {
		t.Errorf("expected script %s, got %s", script, scriptVal)
	}
}

func TestLoadMCPConfig_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	bad := filepath.Join(tmpDir, "bad.json")
	os.WriteFile(bad, []byte("{invalid}"), 0644)

	_, err := loadMCPConfig(bad)
	if err == nil {
		t.Fatal("expected error for invalid json")
	}
	if !strings.Contains(err.Error(), "line") || !strings.Contains(err.Error(), "column") {
		t.Errorf("expected line and column info in error, got %v", err)
	}
}

func TestNeedsTestInteractiveMode(t *testing.T) {
	cases := []struct {
		name string
		opts TestOptions
		want bool
	}{
		{"none", TestOptions{}, true},
		{"all", TestOptions{TestAll: true}, false},
		{"script", TestOptions{ScriptFile: "file"}, false},
		{"config", TestOptions{Config: "cfg"}, false},
	}
	for _, c := range cases {
		if got := needsTestInteractiveMode(&c.opts); got != c.want {
			t.Errorf("%s: expected %v, got %v", c.name, c.want, got)
		}
	}
}

func TestLoadMCPConfig_Valid(t *testing.T) {
	tmpDir := t.TempDir()

	// Direct MCPConfig
	cfg := &core.MCPConfig{Name: "direct", Version: "0.4.1"}
	data, _ := json.Marshal(cfg)
	direct := filepath.Join(tmpDir, "direct.json")
	if err := os.WriteFile(direct, data, 0644); err != nil {
		t.Fatal(err)
	}

	got, err := loadMCPConfig(direct)
	if err != nil {
		t.Fatalf("load direct: %v", err)
	}
	if got.Name != "direct" {
		t.Errorf("expected name 'direct', got %s", got.Name)
	}

	// ProjectConfig
	pc := core.NewProjectConfig()
	pc.Name = "proj"
	pData, _ := json.Marshal(pc)
	proj := filepath.Join(tmpDir, "proj.json")
	if err := os.WriteFile(proj, pData, 0644); err != nil {
		t.Fatal(err)
	}

	got, err = loadMCPConfig(proj)
	if err != nil {
		t.Fatalf("load project: %v", err)
	}
	if got.Name != "proj" {
		t.Errorf("expected name 'proj', got %s", got.Name)
	}
}
