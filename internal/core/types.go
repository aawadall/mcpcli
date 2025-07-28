package core

type Tool struct {
	Name        string
	Description string
}

type Resource struct {
	Name string
	Type string
}

type Capability struct {
	Name    string
	Enabled bool
}

// ResourceType defines allowed types for resources
type ResourceType string

const (
	ResourceTypeDatabase   ResourceType = "database"
	ResourceTypeFilesystem ResourceType = "filesystem"
	ResourceTypeTime       ResourceType = "time"
)

// IsValidResourceType checks if the given type is allowed
func IsValidResourceType(t string) bool {
	switch ResourceType(t) {
	case ResourceTypeDatabase, ResourceTypeFilesystem, ResourceTypeTime:
		return true
	default:
		return false
	}
}

// GetName returns the name of the tool
func (t Tool) GetName() string {
	return t.Name
}

// GetName returns the name of the resource
func (r Resource) GetName() string {
	return r.Name
}

// GetName returns the name of the capability
func (c Capability) GetName() string {
	return c.Name
}