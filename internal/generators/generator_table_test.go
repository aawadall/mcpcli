package generators

import (
	"sort"
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
		{"java", NewJavaGenerator(), "java", []string{"stdio", "rest", "websocket"}},
		{"javascript", NewNodeGenerator(), "javascript", []string{"stdio", "rest", "websocket"}},
		{"python", NewPythonGenerator(), "python", []string{"stdio", "rest", "websocket"}},
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
			sort.Strings(gotTrans)
			expectedTrans := make([]string, len(tt.transports))
			copy(expectedTrans, tt.transports)
			sort.Strings(expectedTrans)
			for i, tr := range expectedTrans {
				if gotTrans[i] != tr {
					t.Errorf("expected transport %s, got %s", tr, gotTrans[i])
				}
			}
		})
	}
}
