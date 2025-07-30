package main

import (
	"os"
	"testing"
)

func TestMainRuns(t *testing.T) {
	// run with --help so command exits normally
	osArgs := []string{"--help"}
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = append([]string{"mcpcli"}, osArgs...)
	main()
}
