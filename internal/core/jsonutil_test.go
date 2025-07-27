package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatJSONError_Syntax(t *testing.T) {
	data := []byte("{\n\"key\": \"value\",,\n}")
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
	ferr := FormatJSONError(data, err, "parse failed")
	if ferr == nil {
		t.Fatal("expected formatted error")
	}
	if !strings.Contains(ferr.Error(), "line") || !strings.Contains(ferr.Error(), "column") {
		t.Errorf("expected line and column info, got %v", ferr)
	}
}
