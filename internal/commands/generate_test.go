package commands

import (
	"os"
	"testing"

	"github.com/aawadall/mcpcli/internal/handlers"
	"github.com/spf13/cobra"
)

func TestNeedsInteractiveMode(t *testing.T) {
	opts := &handlers.GenerateOptions{}
	if !needsInteractiveMode(opts) {
		t.Fatal("expected interactive mode when options empty")
	}
	opts.Name = "name"
	opts.Language = "golang"
	opts.Transport = "stdio"
	if needsInteractiveMode(opts) {
		t.Fatal("expected non-interactive when all options set")
	}
}

func TestPrepareDirectory(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/sub"
	opts := &handlers.GenerateOptions{Name: "proj", Language: "golang", Transport: "stdio", Output: path, Force: false}
	if err := handlers.GenerateProject(opts); err != nil {
		t.Fatalf("GenerateProject failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	opts.Force = true
	if err := handlers.GenerateProject(opts); err != nil {
		t.Fatalf("force GenerateProject failed: %v", err)
	}
}

func TestValidateGenerateOptions_Language(t *testing.T) {
	langs := []string{"golang", "javascript", "java", "python"}
	for _, l := range langs {
		opts := &handlers.GenerateOptions{Name: "p", Language: l, Transport: "stdio"}
		if err := handlers.ValidateGenerateOptions(opts); err != nil {
			t.Fatalf("%s should be valid: %v", l, err)
		}
	}
	opts := &handlers.GenerateOptions{Name: "p", Language: "bad", Transport: "stdio"}
	if err := handlers.ValidateGenerateOptions(opts); err == nil {
		t.Fatal("expected error for unsupported language")
	}
}
func TestAddFlags(t *testing.T) {
	cmd := &cobra.Command{}
	opts := &handlers.GenerateOptions{}

	addFlags(cmd, opts)

	if cmd.Flags().Lookup("name") == nil {
		t.Fatal("expected 'name' flag to be added")
	}
	if cmd.Flags().Lookup("language") == nil {
		t.Fatal("expected 'language' flag to be added")
	}
	if cmd.Flags().Lookup("transport") == nil {
		t.Fatal("expected 'transport' flag to be added")
	}
	if cmd.Flags().Lookup("docker") == nil {
		t.Fatal("expected 'docker' flag to be added")
	}
	if cmd.Flags().Lookup("examples") == nil {
		t.Fatal("expected 'examples' flag to be added")
	}
	if cmd.Flags().Lookup("output") == nil {
		t.Fatal("expected 'output' flag to be added")
	}
	if cmd.Flags().Lookup("force") == nil {
		t.Fatal("expected 'force' flag to be added")
	}

}

func TestNewGenerateCmd(t *testing.T) {
	cmd := NewGenerateCmd()
	if cmd.Use != "generate [name]" {
		t.Fatalf("expected command use to be 'generate [name]', got %s", cmd.Use)
	}
	if len(cmd.Aliases) != 2 || cmd.Aliases[0] != "gen" || cmd.Aliases[1] != "g" {
		t.Fatalf("expected aliases to be 'gen' and 'g', got %v", cmd.Aliases)
	}
	if cmd.Short != "Generate a new MCP server project" {
		t.Fatalf("expected short description to match, got %s", cmd.Short)
	}
}
