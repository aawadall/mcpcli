package core

// Transport defines the transport configuration
type Transport struct {
	Type    string                 `json:"type"`
	Options map[string]interface{} `json:"options,omitempty"`
}
