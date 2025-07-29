package generators

import (
	"testing"
)

// TestGenerators verifies common behavior for all language generators using
// table-driven tests.
func TestGenerators(t *testing.T) {
	tests := []struct {
		name       string
		gen        Generator
		language   string
		transports []string
	}{
		{"go", NewGolangGenerator(), "go", []string{"stdio", "rest", "websocket"}},
		{"java", NewJavaGenerator(), "java", []string{"stdio"}},
		{"javascript", NewNodeGenerator(), "javascript", []string{"stdio"}},
		{"python", NewPythonGenerator(), "python", []string{"stdio"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gen.GetLanguage(); got != tt.language {
				t.Errorf("expected language %s, got %s", tt.language, got)
			}
			gotTrans := tt.gen.GetSupportedTransports()
			if len(gotTrans) != len(tt.transports) {
				t.Fatalf("expected %d transports, got %d", len(tt.transports), len(gotTrans))
			}
			for i, tr := range tt.transports {
				if gotTrans[i] != tr {
					t.Errorf("expected transport %s, got %s", tr, gotTrans[i])
				}
			}
		})
	}
}
