[![Codacy Badge](https://api.codacy.com/project/badge/Grade/25451ebbbd064d70bf905303caefd6d2)](https://app.codacy.com/gh/aawadall/mcpcli?utm_source=github.com&utm_medium=referral&utm_content=aawadall/mcpcli&utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/7b1f68b2c73c49e19e13a7e25f9de2f8)](https://app.codacy.com/gh/aawadall/mcpcli/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/7b1f68b2c73c49e19e13a7e25f9de2f8)](https://app.codacy.com/gh/aawadall/mcpcli/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)
[![Release](https://img.shields.io/github/v/release/aawadall/mcpcli?label=release)](https://github.com/aawadall/mcpcli/releases)

# mcpcli

A CLI tool to scaffold Model Context Protocol (MCP) server projects in Go and other languages. It generates ready-to-use MCP server templates with support for multiple transports, Docker, and example resources/tools.

**Current Version:** v0.4.1 (latest release)

> Learn more about the Model Context Protocol at the [official introduction page](https://modelcontextprotocol.io/introduction).

## Features

- Generate new MCP server projects with a single command
- Supports multiple languages (Go, Node.js, Java, Python)
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
- `internal/`    Core logic, generators, handlers, templates
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

## Testing

Run `go test ./... -cover` to execute the unit tests. Overall coverage should remain above 85%.
Recent additions include tests for `internal/commands/test.go` and error handling cases in `internal/generators/node_test.go` to ensure the CLI testing workflow behaves as expected.
See [the testing guide](doc/testing.md) for more details on running the tests for **mcpcli v0.4.1 (latest)**.
All contributions must maintain this minimum coverage level.

## Contributing

Contributions are welcome! Please open issues or pull requests.

## Releasing

The [release workflow](.github/workflows/release.yml) builds cross-platform archives.
For details on packaging `mcpcli` for Homebrew, APT, and Chocolatey see
[the releasing guide](doc/releasing.md).

## License

MIT License

---

## Further Reading
- [Model Context Protocol Introduction](https://modelcontextprotocol.io/introduction)
