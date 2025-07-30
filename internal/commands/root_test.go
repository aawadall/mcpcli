package commands

import "testing"

// test make root command

var version = "1.0.0" // Example version, replace with actual version if needed
func TestMakeRootCommand(t *testing.T) {
	rootCmd := MakeRootCommand(version)

	if rootCmd.Use != "mcpcli" {
		t.Errorf("expected Use to be 'mcpcli', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short != "A CLI tool for MCP (Model Context Protocol) development" {
		t.Errorf("expected Short description to match, got '%s'", rootCmd.Short)
	}

	if rootCmd.Version != version {
		t.Errorf("expected version to be '%s', got '%s'", version, rootCmd.Version)
	}

	if len(rootCmd.Commands()) == 0 {
		t.Error("expected at least one subcommand")
	}
}

// Test added commands
func TestAddedCommands(t *testing.T) {
	rootCmd := MakeRootCommand(version)
	if len(rootCmd.Commands()) < 2 {
		t.Error("expected at least two subcommands (generate and test)")
	}

	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "generate" || cmd.Name() == "test" {
			if cmd.Use == "" {
				t.Errorf("expected command '%s' to have a valid use description", cmd.Name())
			}
		} else {
			t.Errorf("unexpected command found: %s", cmd.Name())
		}
	}
}

// Test global flags
func TestGlobalFlags(t *testing.T) {
	rootCmd := MakeRootCommand(version)

	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	if verboseFlag == nil || verboseFlag.Name != "verbose" {
		t.Error("expected global flag 'verbose' to be defined")
	}

	quietFlag := rootCmd.PersistentFlags().Lookup("quiet")
	if quietFlag == nil || quietFlag.Name != "quiet" {
		t.Error("expected global flag 'quiet' to be defined")
	}
}

func TestRootCommand_ExecuteUnknown(t *testing.T) {
	rootCmd := MakeRootCommand(version)
	rootCmd.SetArgs([]string{"unknown"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected execute error: %v", err)
	}
}
