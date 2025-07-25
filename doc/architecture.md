# Architecture Overview

## High-Level Design

`mcpcli` is a command-line tool for generating Model Context Protocol (MCP) server projects. It is designed to be extensible, supporting multiple languages and transport methods. The tool uses templates to scaffold new projects with best practices and optional features like Docker and example resources.

## Main Components

- **cmd/**: CLI entrypoint and command definitions (e.g., `generate`)
- **internal/commands/**: Implements CLI commands and their logic
- **internal/core/**: Core types, configuration, and template data structures
- **internal/generators/**: Language-specific project generators
- **internal/templates/**: Project templates for different languages and transports
- **pkg/**: Shared packages (e.g., MCP protocol types)

## Code Generation Flow

1. **User runs `mcpcli generate`** (with flags or interactively)
2. CLI collects options (project name, language, transport, etc.)
3. The appropriate generator (e.g., Go) is selected
4. Templates are rendered with user options and written to the output directory
5. Optional features (Docker, examples) are included as requested

## Generated Project Structure (Go Example)

- `cmd/server/main.go` — Main server entrypoint
- `internal/handlers/` — Request handlers for resources and tools
- `pkg/mcp/` — MCP protocol types and client
- `configs/mcp-config.json` — Server configuration
- `Dockerfile` — Docker support (optional)
- `examples/` — Example usage (optional)

## Generated Project Structure (Node.js Example)

- `src/index.js` — Main server entrypoint
- `src/handlers/` — Request handlers
- `configs/mcp-config.json` — Server configuration
- `Dockerfile` — Docker support (optional)
- `examples/` — Example usage (optional)

## Extensibility

- **Languages**: Add new generators in `internal/generators/` and templates in `internal/templates/`
- **Transports**: Add new transport options in templates and config
- **Features**: Extend CLI flags and template data as needed 