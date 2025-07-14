package core

// TemplateData holds data for template rendering
type TemplateData struct {
	Config      *ProjectConfig
	MCPConfig   *MCPConfig
	PackageName string
	ModuleName  string
	HasDocker   bool
	HasExamples bool
	Timestamp   string
}