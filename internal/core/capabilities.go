package core

// Capabilities defines what the MCP server can do
type Capabilities struct {
	Resources ResourcesCapability `json:"resources,omitempty"`
	Tools     ToolsCapability     `json:"tools,omitempty"`
	Prompts   PromptsCapability   `json:"prompts,omitempty"`
}

// ResourcesCapability defines resource handling capabilities
type ResourcesCapability struct {
	Enabled bool `json:"enabled"`
	Count   int  `json:"count,omitempty"`
}

// ToolsCapability defines tool handling capabilities
type ToolsCapability struct {
	Enabled bool `json:"enabled"`
	Count   int  `json:"count,omitempty"`
}

// PromptsCapability defines prompt handling capabilities
type PromptsCapability struct {
	Enabled bool `json:"enabled"`
	Count   int  `json:"count,omitempty"`
}