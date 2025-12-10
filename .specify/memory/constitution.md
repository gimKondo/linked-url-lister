<!--
  SYNC IMPACT REPORT
  ==================
  Version change: N/A → 1.0.0 (Initial constitution)

  Added Principles:
  - I. Simplicity First
  - II. CLI-First Interface
  - III. Documentation

  Added Sections:
  - Core Principles (3 principles)
  - Technical Standards
  - Development Workflow
  - Governance

  Removed Sections: None (initial creation)

  Templates Status:
  - .specify/templates/plan-template.md: ✅ Compatible (generic Constitution Check section)
  - .specify/templates/spec-template.md: ✅ Compatible (no constitution-specific references)
  - .specify/templates/tasks-template.md: ✅ Compatible (tests marked optional)
  - .specify/templates/checklist-template.md: ✅ Compatible (generic structure)

  Follow-up TODOs: None
-->

# linked-url-lister Constitution

## Core Principles

- ファイル生成やユーザーとのコミュニケーションはすべて日本語で行ってください
- スペック駆動の作業フェイズごと、また実装時にはタスクのフェイズごとにコミットしてください

### I. Simplicity First

Every feature and implementation choice MUST favor simplicity over complexity.

- Start with the simplest solution that could work
- YAGNI (You Aren't Gonna Need It): Do not implement features for hypothetical future requirements
- Avoid premature abstraction: three similar lines of code are preferable to an unnecessary helper function
- External dependencies MUST be justified; prefer standard library when sufficient
- If a design requires extensive documentation to understand, it is too complex

**Rationale**: A CLI tool for URL extraction has a focused scope. Complexity hinders maintainability and increases bugs. Simple code is easier to test, debug, and extend.

### II. CLI-First Interface

The tool MUST provide a clean command-line interface as the primary interaction method.

- Text in/out protocol: read from stdin or file arguments, write results to stdout, errors to stderr
- Support both human-readable and machine-parseable output formats (plain text, JSON)
- Exit codes MUST follow Unix conventions (0 for success, non-zero for errors)
- Flags and arguments MUST follow POSIX conventions where applicable
- Help text (`--help`) MUST be clear and complete

**Rationale**: CLI tools integrate into pipelines, scripts, and automation. A well-designed CLI is the foundation for both human users and programmatic consumption.

### III. Documentation

Code and features MUST be documented sufficiently for users and future maintainers.

- README MUST contain: purpose, installation, basic usage examples, and contribution guidelines
- CLI help text serves as primary user documentation
- Code comments SHOULD explain "why" not "what" when intent is non-obvious
- Breaking changes MUST be documented in release notes

**Rationale**: Good documentation reduces support burden and enables community contributions. For a CLI tool, help text is often the only documentation users consult.

## Technical Standards

**Language**: Go (latest stable version recommended)

**Project Structure**:
- Follow standard Go project layout conventions
- Main package in `cmd/` or root directory
- Internal packages in `internal/` when needed
- Keep package count minimal for a focused tool

**Error Handling**:
- All errors MUST be handled explicitly
- User-facing errors MUST be clear and actionable
- Internal errors SHOULD include context for debugging

**Testing**:
- Tests are encouraged but not mandatory
- When tests exist, they SHOULD cover critical paths and edge cases
- Table-driven tests preferred for Go idiom

## Development Workflow

**Code Quality**:
- Code MUST pass `go fmt` and `go vet` before commit
- Linting with `golangci-lint` is recommended
- Code review required for changes to core parsing logic

**Version Control**:
- Commit messages SHOULD be descriptive and follow conventional format
- Feature branches for non-trivial changes
- Main branch MUST always be in a buildable state

**Releases**:
- Semantic versioning (MAJOR.MINOR.PATCH)
- MAJOR: Breaking CLI interface changes
- MINOR: New features, backward compatible
- PATCH: Bug fixes, documentation updates

## Governance

This constitution defines the non-negotiable standards for the linked-url-lister project.

**Amendment Process**:
1. Proposed amendments MUST be documented with rationale
2. Changes to Core Principles require explicit justification
3. All amendments MUST update the constitution version

**Compliance**:
- Pull requests SHOULD be reviewed against these principles
- Violations of Core Principles require explicit justification and approval
- Use this constitution as the reference for architectural decisions

**Version**: 1.0.0 | **Ratified**: 2025-12-10 | **Last Amended**: 2025-12-10
