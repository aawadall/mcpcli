package generators

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aawadall/mcpcli/internal/core"
)

// BaseTemplateMap returns a mapping of embedded template paths to their output
// locations based on the generator language and configuration.
func BaseTemplateMap(lang string, data *core.TemplateData) (map[string]string, error) {
	switch lang {
	case "go":
		return goTemplateMap(data), nil
	case "javascript":
		return nodeTemplateMap(data), nil
	case "python":
		return pythonTemplateMap(data), nil
	case "java":
		return javaTemplateMap(data), nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}
}

func goTemplateMap(data *core.TemplateData) map[string]string {
	m := map[string]string{
		"templates/go/stdio/go.mod.tmpl":                                              "go.mod",
		fmt.Sprintf("templates/go/%s/cmd/server/main.go.tmpl", data.Config.Transport): filepath.Join("cmd", "server", "main.go"),
		"templates/go/stdio/internal/handlers/mcp.go.tmpl":                            filepath.Join("internal", "handlers", "mcp.go"),
		"templates/go/stdio/internal/resources/filesystem.go.tmpl":                    filepath.Join("internal", "resources", "filesystem.go"),
		"templates/go/stdio/internal/resources/registry.go.tmpl":                      filepath.Join("internal", "resources", "registry.go"),
		"templates/go/stdio/internal/tools/calculator.go.tmpl":                        filepath.Join("internal", "tools", "calculator.go"),
		"templates/go/stdio/pkg/mcp/client.go.tmpl":                                   filepath.Join("pkg", "mcp", "client.go"),
		"templates/go/stdio/pkg/mcp/mcp.go.tmpl":                                      filepath.Join("pkg", "mcp", "mcp.go"),
		"templates/go/stdio/README.md.tmpl":                                           "README.md",
		"templates/go/stdio/configs/mcp-config.json.tmpl":                             filepath.Join("configs", "mcp-config.json"),
		"templates/go/stdio/examples/example.go.tmpl":                                 filepath.Join("examples", "example.go"),
	}
	if data.Config.Docker {
		m["templates/go/stdio/Dockerfile.tmpl"] = "Dockerfile"
		m["templates/go/stdio/dockerignore.tmpl"] = ".dockerignore"
	}
	return m
}

func nodeTemplateMap(data *core.TemplateData) map[string]string {
	m := map[string]string{
		"templates/node/stdio/package.json.tmpl":                                  "package.json",
		fmt.Sprintf("templates/node/%s/src/index.js.tmpl", data.Config.Transport): filepath.Join("src", "index.js"),
		"templates/node/stdio/src/handlers/mcp.js.tmpl":                           filepath.Join("src", "handlers", "mcp.js"),
		"templates/node/stdio/src/resources/registry.js.tmpl":                     filepath.Join("src", "resources", "registry.js"),
		"templates/node/stdio/README.md.tmpl":                                     "README.md",
		"templates/node/stdio/configs/mcp-config.json.tmpl":                       filepath.Join("configs", "mcp-config.json"),
		"templates/node/stdio/examples/example.js.tmpl":                           filepath.Join("examples", "example.js"),
	}
	if data.Config.Docker {
		m["templates/node/stdio/Dockerfile.tmpl"] = "Dockerfile"
		m["templates/node/stdio/dockerignore.tmpl"] = ".dockerignore"
	}
	return m
}

func pythonTemplateMap(data *core.TemplateData) map[string]string {
	m := map[string]string{
		fmt.Sprintf("templates/python/%s/src/main.py.tmpl", data.Config.Transport): filepath.Join("src", "main.py"),
		"templates/python/stdio/src/handlers/mcp.py.tmpl":                          filepath.Join("src", "handlers", "mcp.py"),
		"templates/python/stdio/src/resources/registry.py.tmpl":                    filepath.Join("src", "resources", "registry.py"),
		"templates/python/stdio/README.md.tmpl":                                    "README.md",
		"templates/python/stdio/configs/mcp-config.json.tmpl":                      filepath.Join("configs", "mcp-config.json"),
		"templates/python/stdio/examples/example.py.tmpl":                          filepath.Join("examples", "example.py"),
	}
	if data.Config.Docker {
		m["templates/python/stdio/Dockerfile.tmpl"] = "Dockerfile"
		m["templates/python/stdio/dockerignore.tmpl"] = ".dockerignore"
	}
	return m
}

func javaTemplateMap(data *core.TemplateData) map[string]string {
	pkgPath := filepath.Join(strings.Split(data.PackageName, ".")...)
	m := map[string]string{
		"templates/java/stdio/pom.xml.tmpl":                                                  "pom.xml",
		fmt.Sprintf("templates/java/%s/src/main/java/Main.java.tmpl", data.Config.Transport): filepath.Join("src", "main", "java", pkgPath, "Main.java"),
		"templates/java/stdio/src/main/java/handlers/MCPHandler.java.tmpl":                   filepath.Join("src", "main", "java", pkgPath, "handlers", "MCPHandler.java"),
		"templates/java/stdio/src/main/java/resources/Registry.java.tmpl":                    filepath.Join("src", "main", "java", pkgPath, "resources", "Registry.java"),
		"templates/java/stdio/README.md.tmpl":                                                "README.md",
		"templates/java/stdio/configs/mcp-config.json.tmpl":                                  filepath.Join("configs", "mcp-config.json"),
		"templates/java/stdio/examples/Example.java.tmpl":                                    filepath.Join("examples", "Example.java"),
	}
	if data.Config.Docker {
		m["templates/java/stdio/Dockerfile.tmpl"] = "Dockerfile"
		m["templates/java/stdio/dockerignore.tmpl"] = ".dockerignore"
	}
	return m
}

// ToolTemplate returns the template path for generating a single tool file.
func ToolTemplate(lang string) string {
	switch lang {
	case "go":
		return "templates/go/stdio/internal/tools/tool.go.tmpl"
	case "javascript":
		return "templates/node/stdio/src/tools/tool.js.tmpl"
	case "python":
		return "templates/python/stdio/src/tools/tool.py.tmpl"
	case "java":
		return "templates/java/stdio/src/main/java/tools/Tool.java.tmpl"
	default:
		return ""
	}
}

// ResourceTemplate returns the template path for generating a single resource file.
func ResourceTemplate(lang string) string {
	switch lang {
	case "go":
		return "templates/go/stdio/internal/resources/resource.go.tmpl"
	case "javascript":
		return "templates/node/stdio/src/resources/resource.js.tmpl"
	case "python":
		return "templates/python/stdio/src/resources/resource.py.tmpl"
	case "java":
		return "templates/java/stdio/src/main/java/resources/Resource.java.tmpl"
	default:
		return ""
	}
}

// CapabilityTemplate returns the template path for generating a single capability file.
func CapabilityTemplate(lang string) string {
	switch lang {
	case "go":
		return "templates/go/stdio/internal/capabilities/capability.go.tmpl"
	case "javascript":
		return "templates/node/stdio/src/capabilities/capability.js.tmpl"
	case "python":
		return "templates/python/stdio/src/capabilities/capability.py.tmpl"
	case "java":
		return "templates/java/stdio/src/main/java/capabilities/Capability.java.tmpl"
	default:
		return ""
	}
}
