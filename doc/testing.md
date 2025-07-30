# Running Tests

This guide describes how to run the unit tests for `mcpcli`.

The project targets Go `1.22` and requires no external services. Run all tests with:

```bash
go test ./...
```

To check code coverage, run:

```bash
go test ./internal/core -cover
```
The test suite currently targets **mcpcli v0.4.1 (latest release)** and achieves over 85% coverage for the core package.
Additional tests cover `internal/commands/test.go` and `internal/generators/node_test.go` to validate CLI workflows and Node.js project scaffolding.
Table-driven tests for all generators ensure consistent behavior across languages and maintain over 85% coverage across the project.
