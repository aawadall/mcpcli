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
The test suite currently targets **mcpcli v0.4.2 (latest release)** and achieves over 85% coverage for the core package.
Additional tests cover `internal/commands/test.go` to validate the test command behavior.
Table-driven tests for all generators ensure consistent behavior across languages.
