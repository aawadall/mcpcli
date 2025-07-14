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
