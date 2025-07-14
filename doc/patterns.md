# Design and Implementation Patterns

## 1. Template-Based Code Generation
- Uses Go text/template files for project scaffolding
- Template data is populated from user input and project config
- Supports conditional inclusion (e.g., Docker, examples)

## 2. Generator Interface
- `internal/generators/Generator` interface abstracts project generation for different languages
- Each language implements its own generator (e.g., Go)
- Allows easy extension to new languages

## 3. CLI Command Pattern
- Uses [Cobra](https://github.com/spf13/cobra) for CLI structure
- Each command (e.g., `generate`) is defined as a separate function
- Flags and arguments are managed via Cobra
- Interactive prompts use [survey](https://github.com/AlecAivazis/survey)

## 4. Configuration as Data
- Project and server configuration are represented as Go structs
- Rendered to JSON/YAML for generated projects
- Enables strong typing and validation

## 5. Extensibility
- New features can be added by extending the generator interface and adding new templates
- CLI can be extended with new commands and flags

## 6. Separation of Concerns
- CLI logic, code generation, and template data are separated into different packages
- Templates are kept language- and transport-specific for clarity 