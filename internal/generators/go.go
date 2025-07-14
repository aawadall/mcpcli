package generators

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aawadall/mcpcli/internal/core"
)

// GoGenerator implements the Generator interface for Go projects
type GoGenerator struct {}

func NewGolangGenerator() *GoGenerator {
	return &GoGenerator{}
}

// Generate creates a Go project structure based on the provided configuration
func (g *GoGenerator) Generate(config *core.ProjectConfig) error {
	templateData := config.GetTemplateData()

	// Create directory structure
	if err := g.createDirectoryStructure(config.Output); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Generate files
	if err := g.generateGoMod(config.Output, templateData); err != nil {
		return fmt.Errorf("failed to generate go.mod: %w", err)
	}

	if err := g.generateMainFile(config.Output, templateData); err != nil {
		return fmt.Errorf("failed to generate main file: %w", err)
	}

	if err := g.generateMCPConfig(config.Output, templateData); err != nil {
		return fmt.Errorf("failed to generate MCP config: %w", err)
	}

	if err := g.generatePackageFiles(config.Output, templateData); err != nil {
		return fmt.Errorf("failed to generate package files: %w", err)
	}

	if config.Docker {
		if err := g.generateDockerFiles(config.Output, templateData); err != nil {
			return fmt.Errorf("failed to generate Docker files: %w", err)
		}
	}

	if config.Examples {
		if err := g.generateExamples(config.Output, templateData); err != nil {
			return fmt.Errorf("failed to generate examples: %w", err)
		}
	}

	if err := g.generateREADME(config.Output, templateData); err != nil {
		return fmt.Errorf("failed to generate README: %w", err)
	}

	return nil
}


// GetLanguage returns the language name
func (g *GoGenerator) GetLanguage() string {
	return "go"
}

// GetSupportedTransports returns the list of supported transports for Go
func (g *GoGenerator) GetSupportedTransports() []string {
	return []string{"stdio", "rest", "websocket"}
}

// createDirectoryStructure creates the project directory structure
func (g *GoGenerator) createDirectoryStructure(output string) error {
	dirs := []string{
		"cmd/" + "stdio",
		"cmd/" + "rest",
		"internal/handlers",
		"internal/resources",
		"internal/tools",
		"internal/transport",
		"pkg/mcp",
		"examples",
		"configs",
		"docker",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(output, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

// generateGoMod creates the go.mod file
func (g *GoGenerator) generateGoMod(output string, data *core.TemplateData) error {
	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/sirupsen/logrus v1.9.3
)
`, data.ModuleName)

	return os.WriteFile(filepath.Join(output, "go.mod"), []byte(content), 0644)
}

// generateMainFile creates the main.go file based on transport
func (g *GoGenerator) generateMainFile(output string, data *core.TemplateData) error {
	transportType := data.Config.Transport

	switch transportType {
	case "stdio":
		return g.generateStdioMain(output, data)
	case "rest":
		return g.generateRESTMain(output, data)
	case "websocket":
		return g.generateWebSocketMain(output, data)
	default:
		return fmt.Errorf("unsupported transport: %s", transportType)
	}
}

// generateStdioMain creates the stdio main.go file
func (g *GoGenerator) generateStdioMain(output string, data *core.TemplateData) error {
	content := `package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"{{.ModuleName}}/internal/handlers"
	"{{.ModuleName}}/pkg/mcp"
)

func main() {
	server := mcp.NewServer()
	handler := handlers.NewHandler()

	// Register resources
	server.RegisterResourceHandler(handler.HandleListResources)
	server.RegisterResourceReadHandler(handler.HandleReadResource)

	// Register tools
	server.RegisterToolHandler(handler.HandleListTools)
	server.RegisterCallToolHandler(handler.HandleCallTool)

	// Start stdio server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		
		var request mcp.Request
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			log.Printf("Failed to parse request: %v", err)
			continue
		}

		response := server.HandleRequest(request)
		
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			continue
		}

		fmt.Println(string(responseJSON))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
`

	tmpl, err := template.New("main").Parse(content)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(output, "cmd", "stdio", "main.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateRESTMain creates the REST main.go file
func (g *GoGenerator) generateRESTMain(output string, data *core.TemplateData) error {
	content := `package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"{{.ModuleName}}/internal/handlers"
	"{{.ModuleName}}/pkg/mcp"
)

func main() {
	server := mcp.NewServer()
	handler := handlers.NewHandler()

	// Register handlers
	server.RegisterResourceHandler(handler.HandleListResources)
	server.RegisterResourceReadHandler(handler.HandleReadResource)
	server.RegisterToolHandler(handler.HandleListTools)
	server.RegisterCallToolHandler(handler.HandleCallTool)

	// Setup HTTP routes
	router := mux.NewRouter()
	
	// MCP endpoints
	router.HandleFunc("/mcp/resources", handleListResources(server)).Methods("GET")
	router.HandleFunc("/mcp/resources/{uri}", handleReadResource(server)).Methods("GET")
	router.HandleFunc("/mcp/tools", handleListTools(server)).Methods("GET")
	router.HandleFunc("/mcp/tools/{name}/call", handleCallTool(server)).Methods("POST")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	log.Println("Starting REST server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleListResources(server *mcp.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := mcp.Request{Method: "resources/list"}
		response := server.HandleRequest(request)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleReadResource(server *mcp.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uri := vars["uri"]
		
		request := mcp.Request{
			Method: "resources/read",
			Params: map[string]interface{}{"uri": uri},
		}
		response := server.HandleRequest(request)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleListTools(server *mcp.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := mcp.Request{Method: "tools/list"}
		response := server.HandleRequest(request)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleCallTool(server *mcp.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		toolName := vars["name"]
		
		var params map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		request := mcp.Request{
			Method: "tools/call",
			Params: map[string]interface{}{
				"name":      toolName,
				"arguments": params,
			},
		}
		response := server.HandleRequest(request)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
`

	tmpl, err := template.New("main").Parse(content)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(output, "cmd", "rest", "main.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}


// generateWebSocketMain creates the WebSocket main.go file
func (g *GoGenerator) generateWebSocketMain(output string, data *core.TemplateData) error {
	// WebSocket implementation would go here
	return fmt.Errorf("WebSocket transport not implemented yet")
}

// generateMCPConfig creates the MCP configuration file
func (g *GoGenerator) generateMCPConfig(output string, data *core.TemplateData) error {
	configJSON, err := json.MarshalIndent(data.MCPConfig, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(output, "configs", "mcp-config.json"), configJSON, 0644)
}

// generatePackageFiles creates the package files
func (g *GoGenerator) generatePackageFiles(output string, data *core.TemplateData) error {
	// Generate MCP package
	if err := g.generateMCPPackage(output, data); err != nil {
		return err
	}

	// Generate handlers
	if err := g.generateHandlers(output, data); err != nil {
		return err
	}

	return nil
}


// generateMCPPackage creates the MCP package files
func (g *GoGenerator) generateMCPPackage(output string, data *core.TemplateData) error {
	content := `package mcp

import (
	"encoding/json"
	"fmt"
)

// Request represents an MCP request
type Request struct {
	Method string                 ` + "`json:\"method\"`" + `
	Params map[string]interface{} ` + "`json:\"params,omitempty\"`" + `
	ID     interface{}            ` + "`json:\"id,omitempty\"`" + `
}

// Response represents an MCP response
type Response struct {
	Result interface{} ` + "`json:\"result,omitempty\"`" + `
	Error  *Error      ` + "`json:\"error,omitempty\"`" + `
	ID     interface{} ` + "`json:\"id,omitempty\"`" + `
}

// Error represents an MCP error
type Error struct {
	Code    int    ` + "`json:\"code\"`" + `
	Message string ` + "`json:\"message\"`" + `
}

// Server represents an MCP server
type Server struct {
	resourceHandler     func(Request) Response
	resourceReadHandler func(Request) Response
	toolHandler         func(Request) Response
	callToolHandler     func(Request) Response
}

// NewServer creates a new MCP server
func NewServer() *Server {
	return &Server{}
}

// RegisterResourceHandler registers a resource list handler
func (s *Server) RegisterResourceHandler(handler func(Request) Response) {
	s.resourceHandler = handler
}

/ RegisterResourceReadHandler registers a resource read handler
func (s *Server) RegisterResourceReadHandler(handler func(Request) Response) {
	s.resourceReadHandler = handler
}

// RegisterToolHandler registers a tool list handler
func (s *Server) RegisterToolHandler(handler func(Request) Response) {
	s.toolHandler = handler
}

// RegisterCallToolHandler registers a tool call handler
func (s *Server) RegisterCallToolHandler(handler func(Request) Response) {
	s.callToolHandler = handler
}

// HandleRequest handles an MCP request
func (s *Server) HandleRequest(request Request) Response {
	switch request.Method {
	case "resources/list":
		if s.resourceHandler != nil {
			return s.resourceHandler(request)
		}
	case "resources/read":
		if s.resourceReadHandler != nil {
			return s.resourceReadHandler(request)
		}
	case "tools/list":
		if s.toolHandler != nil {
			return s.toolHandler(request)
		}
	case "tools/call":
		if s.callToolHandler != nil {
			return s.callToolHandler(request)
		}
	default:
		return Response{
			Error: &Error{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", request.Method),
			},
			ID: request.ID,
		}
	}
	return Response{
		Error: &Error{
			Code:    -32600,
			Message: "Invalid Request",
		},
		ID: request.ID,
	}
}
// MarshalJSON marshals the Request to JSON
func (r *Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Method string                 ` + "`json:\"method\"`" + `
		Params map[string]interface{} ` + "`json:\"params,omitempty\"`" + `
		ID     interface{}            ` + "`json:\"id,omitempty\"`" + `
	}{
		Method: r.Method,
		Params: r.Params,
		ID:     r.ID,
	})
}
// MarshalJSON marshals the Response to JSON
func (r *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Result interface{} ` + "`json:\"result,omitempty\"`" + `
		Error  *Error      ` + "`json:\"error,omitempty\"`" + `
		ID     interface{} ` + "`json:\"id,omitempty\"`" + `
	}{
		Result: r.Result,
		Error:  r.Error,
		ID:     r.ID,
	})
}
`

	tmpl, err := template.New("mcp").Parse(content)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(output, "pkg", "mcp", "mcp.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}
// generateHandlers creates the handlers package files
func (g *GoGenerator) generateHandlers(output string, data *core.TemplateData) error {
	content := `package handlers
import (
	"encoding/json"
	"net/http"
	"{{.ModuleName}}/pkg/mcp"
)

// Handler represents the MCP request handler
type Handler struct{}
// NewHandler creates a new MCP handler
func NewHandler() *Handler {
	return &Handler{}
}
// HandleListResources handles the resources list request
func (h *Handler) HandleListResources(req mcp.Request) mcp.Response {
	// Example implementation
	resources := []string{"resource1", "resource2"}
	return mcp.Response{
		Result: resources,
		ID:     req.ID,
	}
}
// HandleReadResource handles the resource read request
func (h *Handler) HandleReadResource(req mcp.Request) mcp.Response {
	// Example implementation
	resourceURI, ok := req.Params["uri"].(string)
	if !ok {
		return mcp.Response{
			Error: &mcp.Error{
				Code:    -32602,
				Message: "Invalid params",
			},
			ID: req.ID,
		}
	}
	resource := map[string]string{"uri": resourceURI, "data": "example data"}
	return mcp.Response{
		Result: resource,
		ID:     req.ID,
	}
}
// HandleListTools handles the tools list request
func (h *Handler) HandleListTools(req mcp.Request) mcp.Response {
	// Example implementation
	tools := []string{"tool1", "tool2"}
	return mcp.Response{
		Result: tools,
		ID:     req.ID,
	}
}
// HandleCallTool handles the tool call request
func (h *Handler) HandleCallTool(req mcp.Request) mcp.Response {
	// Example implementation
	toolName, ok := req.Params["name"].(string)
	if !ok {
		return mcp.Response{
			Error: &mcp.Error{
				Code:    -32602,
				Message: "Invalid params",
			},
			ID: req.ID,
		}
	}
	arguments, ok := req.Params["arguments"].(map[string]interface{})
	if !ok {
		return mcp.Response{
			Error: &mcp.Error{
				Code:    -32602,
				Message: "Invalid arguments",
			},
			ID: req.ID,
		}
	}
	result := map[string]interface{}{
		"tool": toolName,
		"args": arguments,
	}
	return mcp.Response{
		Result: result,
		ID:     req.ID,
	}
}
`

	tmpl, err := template.New("handlers").Parse(content)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(output, "internal", "handlers", "handlers.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}
// generateDockerFiles creates Docker-related files
func (g *GoGenerator) generateDockerFiles(output string, data *core.TemplateData) error {
	dockerfileContent := `FROM golang:1.21
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o mcp-server ./cmd/stdio/main.go
CMD ["./mcp-server"]
`
	dockerfilePath := filepath.Join(output, "docker", "Dockerfile")

	file, err := os.Create(dockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to create Dockerfile: %w", err)
	}
	defer file.Close()
	if _, err := file.WriteString(dockerfileContent); err != nil {
		return fmt.Errorf("failed to write Dockerfile: %w", err)
	}
	dockerignoreContent := `# Ignore Go build artifacts
*.o
*.a
# Ignore binary files
*.out
# Ignore vendor directory
vendor/
# Ignore local configuration files
*.local
# Ignore IDE files
.idea/
# Ignore OS generated files
.DS_Store
# Ignore logs
*.log
# Ignore test files
*.test
# Ignore coverage files
*.cover
# Ignore build directories
build/
# Ignore output directories
dist/
# Ignore temporary files
tmp/
# Ignore cache files
.cache
`
	dockerignorePath := filepath.Join(output, "docker", ".dockerignore")
	file, err = os.Create(dockerignorePath)
	if err != nil {
		return fmt.Errorf("failed to create .dockerignore: %w", err)
	}
	defer file.Close()
	if _, err := file.WriteString(dockerignoreContent); err != nil {
		return fmt.Errorf("failed to write .dockerignore: %w", err)
	}
	return nil
}

// generateExamples creates example files
func (g *GoGenerator) generateExamples(output string, data *core.TemplateData) error {
	exampleContent := `package main
import (
	"fmt"
	"{{.ModuleName}}/pkg/mcp"
	"{{.ModuleName}}/internal/handlers"
)

func main() {
	server := mcp.NewServer()
	handler := handlers.NewHandler()

	// Register resources and tools
	server.RegisterResourceHandler(handler.HandleListResources)
	server.RegisterResourceReadHandler(handler.HandleReadResource)
	server.RegisterToolHandler(handler.HandleListTools)
	server.RegisterCallToolHandler(handler.HandleCallTool)

	// Example request to list resources
	req := mcp.Request{Method: "resources/list"}
	resp := server.HandleRequest(req)
	fmt.Println("Resources:", resp.Result)

	// Example request to read a resource
	req = mcp.Request{Method: "resources/read", Params: map[string]interface{}{"uri": "example/resource"}}
	resp = server.HandleRequest(req)
	fmt.Println("Resource Data:", resp.Result)

	// Example request to list tools
	req = mcp.Request{Method: "tools/list"}
	resp = server.HandleRequest(req)
	fmt.Println("Tools:", resp.Result)

	// Example request to call a tool
	req = mcp.Request{Method: "tools/call", Params: map[string]interface{}{"name": "example_tool", "arguments": map[string]interface{}{"arg1": "value1"}}}
	resp = server.HandleRequest(req)
	fmt.Println("Tool Call Result:", resp.Result)
}
`

	tmpl, err := template.New("example").Parse(exampleContent)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(output, "examples", "example.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}
// generateREADME creates the README.md file
func (g *GoGenerator) generateREADME(output string, data *core.TemplateData) error {
	readmeContent := fmt.Sprintf(`# %s MCP Server
This is a Model Context Protocol (MCP) server implemented in Go.
## Getting Started
To get started, clone the repository and run the following commands:
\`\`\`bash
cd %s
go mod tidy
go run cmd/stdio/main.go
\`\`\`
## Available Commands
- **List Resources**: Send a request to `/mcp/resources` to list all resources.
- **Read Resource**: Send a request to `/mcp/resources/{uri}` to read a specific resource.
- **List Tools**: Send a request to `/mcp/tools` to list all available tools.
- **Call Tool**: Send a POST request to `/mcp/tools/{name}/call` with the tool arguments to call a specific tool.
## Docker Support
To run the server in a Docker container, use the following commands:
\`\`\`bash
docker build -t mcp-server .
docker run -p 8080:8080 mcp-server
\`\`\`
## Examples
You can find example usage in the `examples` directory.
## License
This project is licensed under the MIT License.
`, data.Config.Name, data.Config.Output)
	readmePath := filepath.Join(output, "README.md")
	file, err := os.Create(readmePath)
	if err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}
	defer file.Close()
	if _, err := file.WriteString(readmeContent); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}
	return nil
}
// sanitizePackageName ensures the package name is valid for Go
func sanitizePackageName(name string) string {
	// Replace invalid characters with underscores
	sanitized := ""
	for _, r := range name {
		if r == '-' || r == ' ' {
			sanitized += "_"
		} else if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			sanitized += string(r)
		} else {
			sanitized += "_"
		}
	}
	return sanitized
}

// GetPackageName returns a sanitized package name for Go
func (g *GoGenerator) GetPackageName(name string) string {
	return sanitizePackageName(name)
}

// GetModuleName returns the module name for Go
func (g *GoGenerator) GetModuleName(name string) string {
	// For Go, the module name is typically the same as the package name
	return sanitizePackageName(name)
}

// GetTimestamp returns the current timestamp in RFC3339 format
func (g *GoGenerator) GetTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// GetDefaultTransportOptions returns default options for the specified transport type
func (g *GoGenerator) GetDefaultTransportOptions(transportType string) map[string]interface{} {
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

// GetDefaultDockerfileContent returns the default Dockerfile content for Go projects
func (g *GoGenerator) GetDefaultDockerfileContent() string {
	return `FROM golang:1.21
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o mcp-server ./cmd/stdio/main.go
CMD ["./mcp-server"]
`
}
// GetDefaultDockerignoreContent returns the default .dockerignore content for Go projects
func (g *GoGenerator) GetDefaultDockerignoreContent() string {
	return `# Ignore Go build artifacts
*.o
*.a
# Ignore binary files
*.out
# Ignore vendor directory
vendor/
# Ignore local configuration files
*.local
# Ignore IDE files
.idea/
# Ignore OS generated files
.DS_Store
# Ignore logs
*.log
# Ignore test files
*.test
# Ignore coverage files
*.cover
# Ignore build directories
build/
# Ignore output directories
dist/
# Ignore temporary files
tmp/
# Ignore cache files
.cache
`
}
// GetDefaultExampleContent returns the default example content for Go projects
func (g *GoGenerator) GetDefaultExampleContent() string {
	return `package main
import (
	"fmt"
	"{{.ModuleName}}/pkg/mcp"
	"{{.ModuleName}}/internal/handlers"
)
func main() {
	server := mcp.NewServer()
	handler := handlers.NewHandler()

	// Register resources and tools
	server.RegisterResourceHandler(handler.HandleListResources)
	server.RegisterResourceReadHandler(handler.HandleReadResource)
	server.RegisterToolHandler(handler.HandleListTools)
	server.RegisterCallToolHandler(handler.HandleCallTool)

	// Example request to list resources
	req := mcp.Request{Method: "resources/list"}
	resp := server.HandleRequest(req)
	fmt.Println("Resources:", resp.Result)

	// Example request to read a resource
	req = mcp.Request{Method: "resources/read", Params: map[string]interface{}{"uri": "example/resource"}}
	resp = server.HandleRequest(req)
	fmt.Println("Resource Data:", resp.Result)

	// Example request to list tools
	req = mcp.Request{Method: "tools/list"}
	resp = server.HandleRequest(req)
	fmt.Println("Tools:", resp.Result)

	// Example request to call a tool
	req = mcp.Request{Method: "tools/call", Params: map[string]interface{}{"name": "example_tool", "arguments": map[string]interface{}{"arg1": "value1"}}}
	resp = server.HandleRequest(req)
	fmt.Println("Tool Call Result:", resp.Result)
}
`
}
