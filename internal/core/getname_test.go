package core

import "testing"

func TestToolGetName(t *testing.T) {
	tool := Tool{Name: "mytool"}
	if tool.GetName() != "mytool" {
		t.Errorf("expected 'mytool', got %s", tool.GetName())
	}
}

func TestResourceGetName(t *testing.T) {
	r := Resource{Name: "db"}
	if r.GetName() != "db" {
		t.Errorf("expected 'db', got %s", r.GetName())
	}
}

func TestCapabilityGetName(t *testing.T) {
	c := Capability{Name: "cap", Enabled: true}
	if c.GetName() != "cap" {
		t.Errorf("expected 'cap', got %s", c.GetName())
	}
}
