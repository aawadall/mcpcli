package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/aawadall/mcpcli/internal/core"
)

func writeTempConfig(t *testing.T, dir string) string {
	cfg := &core.MCPConfig{
		Name:      "test",
		Version:   "0.4.0",
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
