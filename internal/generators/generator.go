package generators

import "github.com/aawadall/mcpcli/internal/core"

// Generator interface for different language generators
type Generator interface {
	Generate(config *core.ProjectConfig) error
	GetLanguage() string
	GetSupportedTransports() []string
}
