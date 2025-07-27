package core

import "testing"

func TestIsValidResourceType(t *testing.T) {
	valid := []string{"database", "filesystem", "time"}
	for _, v := range valid {
		if !IsValidResourceType(v) {
			t.Errorf("expected %s to be valid", v)
		}
	}
	if IsValidResourceType("invalid") {
		t.Error("unexpected valid result for invalid type")
	}
}
