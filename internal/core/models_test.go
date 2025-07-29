package core

import (
	"encoding/json"
	"testing"
)

func TestNewMCPConfigDefaults(t *testing.T) {
	cfg := NewMCPConfig("svc", "1.0.0", "desc", nil, nil)
	if cfg.Schema == "" || cfg.Schema != "https://schemas.modelcontextprotocol.org/server-config.json" {
		t.Errorf("unexpected schema %s", cfg.Schema)
	}
	if cfg.Name != "svc" || cfg.Version != "1.0.0" || cfg.Description != "desc" {
		t.Error("parameters not propagated")
	}
	if cfg.License != "MIT" {
		t.Error("license should default to MIT")
	}
	if cfg.Transport.Type != "stdio" {
		t.Errorf("transport default not stdio: %s", cfg.Transport.Type)
	}
	if !cfg.Capabilities.Resources.Enabled || !cfg.Capabilities.Tools.Enabled || !cfg.Capabilities.Prompts.Enabled {
		t.Error("capabilities not enabled by default")
	}
}

func TestMCPConfigSetTransport(t *testing.T) {
	cfg := NewMCPConfig("svc", "1.0.0", "desc", nil, nil)
	opts := map[string]interface{}{"port": 8080}
	cfg.SetTransport("rest", opts)
	if cfg.Transport.Type != "rest" {
		t.Errorf("expected rest transport, got %s", cfg.Transport.Type)
	}
	if v := cfg.Transport.Options["port"]; v != 8080 {
		t.Errorf("expected port option 8080, got %v", v)
	}
}

func TestRequestResponseJSON(t *testing.T) {
	req := Request{Method: "ping", Params: map[string]interface{}{"msg": "hi"}, ID: 1}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	var got Request
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal request: %v", err)
	}
	if got.Method != req.Method {
		t.Errorf("method mismatch: %+v", got)
	}
	if id, ok := got.ID.(float64); !ok || id != 1 {
		t.Errorf("unexpected ID: %#v", got.ID)
	}
	resp := Response{Result: "pong", ID: 1}
	rdata, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}
	var rgot Response
	if err := json.Unmarshal(rdata, &rgot); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if rgot.Result != resp.Result {
		t.Errorf("result mismatch: %+v", rgot)
	}
	if id, ok := rgot.ID.(float64); !ok || id != 1 {
		t.Errorf("unexpected ID: %#v", rgot.ID)
	}
}

func TestTransportMarshal(t *testing.T) {
	tr := Transport{Type: "rest", Options: map[string]interface{}{"host": "localhost"}}
	data, err := json.Marshal(tr)
	if err != nil {
		t.Fatalf("marshal transport: %v", err)
	}
	var out Transport
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal transport: %v", err)
	}
	if out.Type != tr.Type || out.Options["host"] != "localhost" {
		t.Errorf("unexpected transport round trip: %+v", out)
	}
}

func TestCapabilitiesMarshal(t *testing.T) {
	caps := Capabilities{Resources: ResourcesCapability{Enabled: true, Count: 1}}
	b, err := json.Marshal(caps)
	if err != nil {
		t.Fatalf("marshal capabilities: %v", err)
	}
	var out Capabilities
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal capabilities: %v", err)
	}
	if !out.Resources.Enabled || out.Resources.Count != 1 {
		t.Errorf("unexpected capabilities: %+v", out)
	}
}
