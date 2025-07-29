package commands

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestNeedsInteractiveMode(t *testing.T) {
	opts := &GenerateOptions{}
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

func TestContains(t *testing.T) {
	items := []string{"a", "b"}
	if !contains(items, "a") {
		t.Fatal("expected to find item")
	}
	if contains(items, "c") {
		t.Fatal("did not expect to find item")
	}
}

func TestPrepareDirectory(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/sub"
	if err := prepareDirectory(path, false); err != nil {
		t.Fatalf("prepareDirectory failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	// create again without force should error
	if err := prepareDirectory(path, false); err == nil {
		t.Fatal("expected error when directory exists")
	}
	if err := prepareDirectory(path, true); err != nil {
		t.Fatalf("prepareDirectory force failed: %v", err)
	}
}

func TestSelectGenerator(t *testing.T) {
	langs := []string{"golang", "javascript", "java", "python"}
	for _, l := range langs {
		if gen, err := selectGenerator(l); err != nil || gen == nil {
			t.Fatalf("generator for %s not returned", l)
		}
	}
	if _, err := selectGenerator("bad"); err == nil {
		t.Fatal("expected error for unsupported language")
	}
}
func TestAddFlags(t *testing.T) {
	cmd := &cobra.Command{}
	opts := &GenerateOptions{}

	
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
	// END: Test adding flags
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

