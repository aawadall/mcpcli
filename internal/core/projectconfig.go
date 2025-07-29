package core

import "time"

type ProjectConfig struct {
	Name        string    `json:"name"`
	Language    string    `json:"language"`
	Transport   string    `json:"transport"`
	Docker      bool      `json:"docker"`
	Examples    bool      `json:"examples"`
	Output      string    `json:"output"`
	Author      string    `json:"author,omitempty"`
	Description string    `json:"description,omitempty"`
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"created_at"`

	Tools        []Tool       `json:"tools,omitempty"`
	Resources    []Resource   `json:"resources,omitempty"`
	Capabilities []Capability `json:"capabilities,omitempty"`
}

// NewProjectConfig creates a new project configuration with defaults
func NewProjectConfig() *ProjectConfig {
	return &ProjectConfig{
		Version:   "0.4.2",
		CreatedAt: time.Now(),
	}
}

// GetTemplateData creates template data from the project config
func (pc *ProjectConfig) GetTemplateData() *TemplateData {
	mcpConfig := NewMCPConfig(pc.Name, pc.Version, pc.Description, pc.Tools, pc.Resources)
	mcpConfig.SetTransport(pc.Transport, getTransportOptions(pc.Transport))

	return &TemplateData{
		Config:      pc,
		MCPConfig:   mcpConfig,
		PackageName: sanitizePackageName(pc.Name),
		ModuleName:  pc.Name,
		HasDocker:   pc.Docker,
		HasExamples: pc.Examples,
		Timestamp:   pc.CreatedAt.Format(time.RFC3339),
	}
}

// getTransportOptions returns default options for each transport type
func getTransportOptions(transportType string) map[string]interface{} {
	switch transportType {
	case "rest":
		return map[string]interface{}{
			"port": 8080,
			"host": "localhost",
		}
	case "websocket":
		return map[string]interface{}{
			"port": 8081,
			"path": "/ws",
		}
	default:
		return nil
	}
}

// sanitizePackageName ensures the package name is valid for Go
func sanitizePackageName(name string) string {
	// Simple sanitization - replace hyphens with underscores
	// and ensure it starts with a letter
	result := ""
	for i, r := range name {
		if i == 0 {
			if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
				result += string(r)
			} else {
				result += "pkg"
			}
		} else {
			if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
				result += string(r)
			} else if r == '-' || r == '_' {
				result += "_"
			}
		}
	}
	return result
}
