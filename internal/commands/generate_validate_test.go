package commands

import (
	"github.com/aawadall/mcpcli/internal/handlers"
	"testing"
)

func TestValidateOptions(t *testing.T) {
	opts := &handlers.GenerateOptions{Name: "proj", Language: "golang", Transport: "stdio"}
	if err := handlers.ValidateGenerateOptions(opts); err != nil {
		t.Fatalf("valid options returned error: %v", err)
	}

	opts.Language = "invalid"
	if err := handlers.ValidateGenerateOptions(opts); err == nil {
		t.Fatal("expected error for invalid language")
	}

	opts.Language = "golang"
	opts.Transport = "bad"
	if err := handlers.ValidateGenerateOptions(opts); err == nil {
		t.Fatal("expected error for invalid transport")
	}
}
