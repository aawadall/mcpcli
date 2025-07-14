# mcpcli

A CLI tool to scaffold Model Context Protocol (MCP) server projects in Go and other languages. It generates ready-to-use MCP server templates with support for multiple transports, Docker, and example resources/tools.

## Features
- Generate new MCP server projects with a single command
- Supports multiple languages (currently: Go)
- Choose transport method (currently: stdio)
- Optional Docker support
- Example resources and tools included
- Interactive and non-interactive modes

## Installation

Clone the repository and build the CLI:

```bash
git clone https://github.com/aawadall/mcpcli.git
cd mcpcli
go build -o mcpcli ./cmd/mcpcli
```

## Usage

Generate a new MCP server project:

```bash
./mcpcli generate my-server --language golang --transport stdio --docker --examples
```

Or use interactive mode:

```bash
./mcpcli generate
```

### Flags
- `--name, -n`         Project name
- `--language, -l`     Programming language (e.g., golang)
- `--transport, -t`    Transport method (e.g., stdio)
- `--docker, -d`       Include Docker support
- `--examples, -e`     Include example resources and tools
- `--output, -o`       Output directory
- `--force, -f`        Overwrite existing directory

## Project Structure
- `cmd/`         Entrypoint for the CLI
- `internal/`    Core logic, generators, templates
- `pkg/`         Shared packages

## Generated Project Example

A generated Go MCP server project includes:
- `cmd/server/main.go` - Main server entrypoint
- `internal/handlers/` - Request handlers
- `pkg/mcp/`           - MCP protocol types and client
- `configs/mcp-config.json` - Server configuration
- `Dockerfile`         - Docker support (optional)
- `examples/`          - Example usage (optional)

## Contributing

Contributions are welcome! Please open issues or pull requests.

## License

MIT License 