package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateCmd_FlagParsingAndExecution(t *testing.T) {
	tmpDir := t.TempDir()
	output := filepath.Join(tmpDir, "out")

	cmd := NewGenerateCmd()
	cmd.SetArgs([]string{
		"--name", "myproj",
		"--language", "golang",
		"--transport", "stdio",
		"--output", output,
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	name, _ := cmd.Flags().GetString("name")
	if name != "myproj" {
		t.Errorf("expected name 'myproj', got %s", name)
	}
	lang, _ := cmd.Flags().GetString("language")
	if lang != "golang" {
		t.Errorf("expected language 'golang', got %s", lang)
	}

	if stat, err := os.Stat(output); err != nil || !stat.IsDir() {
		t.Errorf("expected output directory to be created: %v", err)
	}
}
