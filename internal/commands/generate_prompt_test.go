package commands

import (
	"testing"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/aawadall/mcpcli/internal/core"
	"github.com/aawadall/mcpcli/internal/handlers"
)

func TestBuildAndApplyBasicQuestions(t *testing.T) {
	opts := &handlers.GenerateOptions{}
	qs := buildBasicQuestions(opts)
	if len(qs) == 0 {
		t.Fatal("expected questions for empty options")
	}
	answers := basicAnswers{
		Name:      "myproj",
		Language:  "golang",
		Transport: "stdio",
		Docker:    true,
		Examples:  true,
		Output:    "out",
	}
	applyBasicAnswers(opts, answers)
	if opts.Name != "myproj" || opts.Language != "golang" || opts.Output != "out" {
		t.Fatalf("answers not applied: %+v", opts)
	}
}

func TestPromptForTools(t *testing.T) {
	origAskOne := survey.AskOne
	origAsk := survey.Ask
	defer func() { survey.AskOne = origAskOne; survey.Ask = origAsk }()
	call := 0
	survey.AskOne = func(p interface{}, r interface{}, _ ...interface{}) error {
		if b, ok := r.(*bool); ok {
			if call == 0 {
				*b = true // add one tool
			} else {
				*b = false // stop loop
			}
			call++
		}
		return nil
	}
	survey.Ask = func(qs interface{}, resp interface{}, _ ...interface{}) error {
		if tool, ok := resp.(*core.Tool); ok {
			tool.Name = "t"
			tool.Description = "d"
		}
		return nil
	}
	opts := &handlers.GenerateOptions{}
	if err := promptForTools(opts); err != nil {
		t.Fatalf("promptForTools error: %v", err)
	}
	if len(opts.Tools) != 1 || opts.Tools[0].Name != "t" {
		t.Fatalf("tool not added: %+v", opts.Tools)
	}
}

func TestPromptForResourcesAndCapabilities(t *testing.T) {
	origOne := survey.AskOne
	origAsk := survey.Ask
	defer func() { survey.AskOne = origOne; survey.Ask = origAsk }()
	call := 0
	survey.AskOne = func(p interface{}, r interface{}, _ ...interface{}) error {
		if b, ok := r.(*bool); ok {
			if call == 0 {
				*b = true // add one item
			} else {
				*b = false // stop loops
			}
			call++
		}
		return nil
	}
	survey.Ask = func(qs interface{}, resp interface{}, _ ...interface{}) error {
		switch v := resp.(type) {
		case *core.Resource:
			v.Name = "r"
			v.Type = string(core.ResourceTypeFilesystem)
		case *core.Capability:
			v.Name = "c"
			v.Enabled = true
		}
		return nil
	}
	opts := &handlers.GenerateOptions{}
	if err := promptForResources(opts); err != nil {
		t.Fatalf("promptForResources error: %v", err)
	}
	// reset call counter for capabilities
	call = 0
	if err := promptForCapabilities(opts); err != nil {
		t.Fatalf("promptForCapabilities error: %v", err)
	}
	if len(opts.Resources) != 1 || opts.Resources[0].Name != "r" {
		t.Fatalf("resource not added: %+v", opts.Resources)
	}
	if len(opts.Capabilities) != 1 || opts.Capabilities[0].Name != "c" {
		t.Fatalf("capability not added: %+v", opts.Capabilities)
	}
}
func TestPromptForOptions(t *testing.T) {
	origAskOne := survey.AskOne
	origAsk := survey.Ask
	defer func() { survey.AskOne = origAskOne; survey.Ask = origAsk }()
	survey.AskOne = func(p interface{}, r interface{}, _ ...interface{}) error {
		if b, ok := r.(*bool); ok {
			if c, ok := p.(*survey.Confirm); ok {
				switch c.Message {
				case "Would you like to add tools?", "Would you like to add resources?", "Would you like to add capabilities?":
					*b = true
				default:
					*b = false
				}
			}
		} else if s, ok := r.(*string); ok {
			*s = "ans"
		}
		return nil
	}
	survey.Ask = func(qs interface{}, resp interface{}, _ ...interface{}) error {
		switch v := resp.(type) {
		case *basicAnswers:
			v.Name = "p"
			v.Language = "golang"
			v.Transport = "stdio"
			v.Docker = true
			v.Examples = true
			v.Output = "out"
		case *core.Tool:
			v.Name = "t"
		case *core.Resource:
			v.Name = "r"
			v.Type = string(core.ResourceTypeFilesystem)
		case *core.Capability:
			v.Name = "c"
		}
		return nil
	}
	opts := &handlers.GenerateOptions{}
	if err := promptForOptions(opts); err != nil {
		t.Fatalf("promptForOptions err: %v", err)
	}
	if opts.Name == "" || len(opts.Tools) == 0 || len(opts.Resources) == 0 || len(opts.Capabilities) == 0 {
		t.Fatalf("options not fully populated: %+v", opts)
	}
}
