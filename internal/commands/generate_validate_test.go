package commands

import "testing"

func TestValidateOptions(t *testing.T) {
	opts := &GenerateOptions{Name: "proj", Language: "golang", Transport: "stdio"}
	if err := validateOptions(opts); err != nil {
		t.Fatalf("valid options returned error: %v", err)
	}

	opts.Language = "invalid"
	if err := validateOptions(opts); err == nil {
		t.Fatal("expected error for invalid language")
	}

	opts.Language = "golang"
	opts.Transport = "bad"
	if err := validateOptions(opts); err == nil {
		t.Fatal("expected error for invalid transport")
	}
}
