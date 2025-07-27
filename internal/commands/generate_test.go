package commands

import (
	"os"
	"testing"
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
