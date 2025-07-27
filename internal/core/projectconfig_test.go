package core

import "testing"

func TestSanitizePackageName(t *testing.T) {
	cases := map[string]string{
		"my-project": "my_project",
		"1bad-name":  "pkgbad_name",
		"HelloWorld": "HelloWorld",
	}
	for input, expected := range cases {
		if out := sanitizePackageName(input); out != expected {
			t.Errorf("%s => %s, want %s", input, out, expected)
		}
	}
}

func TestGetTransportOptions(t *testing.T) {
	rest := getTransportOptions("rest")
	if rest["port"] != 8080 || rest["host"] != "localhost" {
		t.Errorf("unexpected rest options: %v", rest)
	}
	ws := getTransportOptions("websocket")
	if ws["port"] != 8081 || ws["path"] != "/ws" {
		t.Errorf("unexpected websocket options: %v", ws)
	}
	if getTransportOptions("stdio") != nil {
		t.Error("expected nil for stdio transport")
	}
}

func TestNewProjectConfig_GetTemplateData(t *testing.T) {
	pc := NewProjectConfig()
	pc.Name = "example"
	pc.Transport = "rest"
	data := pc.GetTemplateData()
	if data.Config != pc {
		t.Error("template data should reference original config")
	}
	if data.PackageName != "example" {
		t.Errorf("expected package name 'example', got %s", data.PackageName)
	}
	if data.MCPConfig.Transport.Type != "rest" {
		t.Errorf("expected transport 'rest', got %s", data.MCPConfig.Transport.Type)
	}
	opts := data.MCPConfig.Transport.Options
	if opts["port"] != 8080 {
		t.Errorf("unexpected port %v", opts["port"])
	}
}
