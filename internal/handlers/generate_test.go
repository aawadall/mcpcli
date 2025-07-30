package handlers

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateGenerateOptions(t *testing.T) {
	opts := &GenerateOptions{Name: "proj", Language: "golang", Transport: "stdio"}
	if err := ValidateGenerateOptions(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	opts.Language = "unknown"
	if err := ValidateGenerateOptions(opts); err == nil {
		t.Fatal("expected error for invalid language")
	}
}

func TestGenerateProjectCreatesDir(t *testing.T) {
	tmp := t.TempDir()
	out := filepath.Join(tmp, "proj")
	opts := &GenerateOptions{Name: "proj", Language: "golang", Transport: "stdio", Output: out}
	if err := GenerateProject(opts); err != nil {
		t.Fatalf("GenerateProject failed: %v", err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("expected directory %s to exist", out)
	}
}

func TestPrepareDirectoryForce(t *testing.T) {
	tmp := t.TempDir()
	// create subdir
	path := tmp + "/sub"
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
	// should fail without force
	if err := prepareDirectory(path, false); err == nil {
		t.Fatal("expected error when directory exists without force")
	}
	// should succeed with force
	if err := prepareDirectory(path, true); err != nil {
		t.Fatalf("force failed: %v", err)
	}
}

func TestSelectGeneratorUnsupported(t *testing.T) {
	if _, err := selectGenerator("badlang"); err == nil {
		t.Fatal("expected error for unsupported language")
	}
}
func TestSelectGeneratorSupported(t *testing.T) {
	langs := []string{"golang", "javascript", "java", "python"}
	for _, l := range langs {
		if _, err := selectGenerator(l); err != nil {
			t.Fatalf("generator for %s not found: %v", l, err)
		}
	}
}

func TestPrintNextSteps(t *testing.T) {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()
	printNextSteps(&GenerateOptions{Name: "p", Language: "golang", Output: "out"})
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if !strings.Contains(buf.String(), "Next steps") {
		t.Fatalf("output missing next steps")
	}
}

func captureGenOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()
	f()
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintNextSteps_OtherLanguages(t *testing.T) {
	cases := []struct {
		lang   string
		expect string
	}{
		{"javascript", "npm install"},
		{"java", "mvn package"},
		{"python", "python src/main.py"},
	}
	for _, c := range cases {
		out := captureGenOutput(func() {
			printNextSteps(&GenerateOptions{Name: "p", Language: c.lang, Output: "out"})
		})
		if !strings.Contains(out, c.expect) {
			t.Fatalf("expected %s instructions", c.lang)
		}
	}
}
