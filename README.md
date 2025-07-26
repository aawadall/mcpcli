# mcpcli

A CLI tool to scaffold Model Context Protocol (MCP) server projects in Go and other languages. It generates ready-to-use MCP server templates with support for multiple transports, Docker, and example resources/tools.

> Learn more about the Model Context Protocol at the [official introduction page](https://modelcontextprotocol.io/introduction).

## Features
- Generate new MCP server projects with a single command
- Supports multiple languages (Go, Node.js, Java); Python support planned
- Choose transport method (stdio, rest, websocket)
- Optional Docker support
- Example resources and tools included
- Interactive and non-interactive modes
- Test MCP server resources, tools, and capabilities

## Installation

Clone the repository and build the CLI:

```bash
git clone https://github.com/aawadall/mcpcli.git
cd mcpcli
go build -o mcpcli ./cmd/mcpcli
```

## Available Commands

- `generate` (aliases: `gen`, `g`): Generate a new MCP server project
- `test`: Test MCP server resources, tools, capabilities, and initialization

## Usage

### Generate a new MCP server project

```bash
./mcpcli generate my-server --language golang --transport stdio --docker --examples
```

Or use interactive mode (if required options are missing):

```bash
./mcpcli generate
```

#### Generate Flags
- `--name, -n`         Project name
- `--language, -l`     Programming language (`golang`, `python`, `java`, `javascript`/Node.js)
- `--transport, -t`    Transport method (`stdio`, `rest`, `websocket`)
- `--docker, -d`       Include Docker support
- `--examples, -e`     Include example resources and tools
- `--output, -o`       Output directory (default: project name)
- `--force, -f`        Overwrite existing directory

### Test an MCP server

```bash
./mcpcli test --config configs/mcp-config.json --all
```

Or use interactive mode (if no flags are provided):

```bash
./mcpcli test
```

#### Test Flags
- `--config, -c`         Path to MCP configuration file
- `--all`                Test all components (resources, tools, capabilities, init)
- `--resources`          Test resources
- `--tools`              Test tools
- `--capabilities`       Test capabilities
- `--init`               Test initialization
- `--script, -f`         Path to test script file

### Global Flags
- `--verbose, -v`   Enable verbose output
- `--quiet, -q`     Suppress output

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

A generated Node.js MCP server project includes:
- `src/index.js` - Main server entrypoint
- `src/handlers/` - Request handlers
- `configs/mcp-config.json` - Server configuration
- `Dockerfile` - Docker support (optional)

## Contributing

Contributions are welcome! Please open issues or pull requests.

## License

MIT License

---

## Further Reading
- [Model Context Protocol Introduction](https://modelcontextprotocol.io/introduction) 