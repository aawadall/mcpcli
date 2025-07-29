package core

// TemplateData holds information used when rendering server templates
type TemplateData struct {
	Config      *ProjectConfig
	MCPConfig   *MCPConfig
	PackageName string
	ModuleName  string
	HasDocker   bool
	HasExamples bool
	Timestamp   string
}
