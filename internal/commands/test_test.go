package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/aawadall/mcpcli/internal/core"
	"github.com/aawadall/mcpcli/internal/handlers"
)

func writeTempConfig(t *testing.T, dir string) string {
	cfg := &core.MCPConfig{
		Name:      "test",
		Version:   "0.4.2",
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

	_, err := handlers.LoadMCPConfig(bad)
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
		opts handlers.TestOptions
		want bool
	}{
		{"none", handlers.TestOptions{}, true},
		{"all", handlers.TestOptions{TestAll: true}, false},
		{"script", handlers.TestOptions{ScriptFile: "file"}, false},
		{"config", handlers.TestOptions{Config: "cfg"}, false},
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
	cfg := &core.MCPConfig{Name: "direct", Version: "0.4.2"}
	data, _ := json.Marshal(cfg)
	direct := filepath.Join(tmpDir, "direct.json")
	if err := os.WriteFile(direct, data, 0644); err != nil {
		t.Fatal(err)
	}

	got, err := handlers.LoadMCPConfig(direct)
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

	got, err = handlers.LoadMCPConfig(proj)
	if err != nil {
		t.Fatalf("load project: %v", err)
	}
	if got.Name != "proj" {
		t.Errorf("expected name 'proj', got %s", got.Name)
	}
}

func TestLoadMCPConfig_FileNotFound(t *testing.T) {
	_, err := handlers.LoadMCPConfig("no-such-file.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewTestCmd_HasFlags(t *testing.T) {
	cmd := NewTestCmd()
	flags := []string{"config", "all", "resources", "tools", "capabilities", "init", "script"}
	for _, f := range flags {
		if cmd.Flags().Lookup(f) == nil {
			t.Errorf("flag %s not defined", f)
		}
	}
}
func TestPromptForTestOptions_SelectAll(t *testing.T) {
	orig := survey.AskOne
	defer func() { survey.AskOne = orig }()
	call := 0
	survey.AskOne = func(prompt interface{}, response interface{}, opts ...interface{}) error {
		call++
		switch call {
		case 1:
			if r, ok := response.(*[]string); ok {
				*r = []string{"All"}
			}
		case 2:
			if r, ok := response.(*string); ok {
				*r = "mycfg.json"
			}
		}
		return nil
	}
	opts := &handlers.TestOptions{}
	if err := promptForTestOptions(opts); err != nil {
		t.Fatalf("prompt failed: %v", err)
	}
	if !opts.TestAll {
		t.Errorf("expected TestAll true")
	}
	if opts.Config != "mycfg.json" {
		t.Errorf("expected config 'mycfg.json', got %s", opts.Config)
	}
}
