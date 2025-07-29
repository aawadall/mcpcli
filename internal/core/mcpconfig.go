package core

// MCPConfig represents the configuration file used by an MCP server.
type MCPConfig struct {
	Schema       string       `json:"$schema"`
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Description  string       `json:"description"`
	Author       string       `json:"author,omitempty"`
	License      string       `json:"license,omitempty"`
	Repository   string       `json:"repository,omitempty"`
	Transport    Transport    `json:"transport"`
	Capabilities Capabilities `json:"capabilities"`
	Tools        []Tool       `json:"tools,omitempty"`
	Resources    []Resource   `json:"resources,omitempty"`
}

// NewMCPConfig returns an MCPConfig initialized with default values
// including a MIT license and stdio transport.
func NewMCPConfig(name, version, description string, tools []Tool, resources []Resource) *MCPConfig {
	return &MCPConfig{
		Schema:      "https://schemas.modelcontextprotocol.org/server-config.json",
		Name:        name,
		Version:     version,
		Description: description,
		License:     "MIT",
		Transport: Transport{
			Type: "stdio",
		},
		Capabilities: Capabilities{
			Resources: ResourcesCapability{Enabled: true},
			Tools:     ToolsCapability{Enabled: true},
			Prompts:   PromptsCapability{Enabled: true},
		},
		Tools:     tools,
		Resources: resources,
	}
}

// SetTransport updates the transport configuration
func (c *MCPConfig) SetTransport(transportType string, options map[string]interface{}) {
	c.Transport.Type = transportType
	c.Transport.Options = options
}
