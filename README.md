# Krewire TUI (`github.com/krewire/tui`)

Lightweight, pure-Go Terminal User Interface (TUI) framework and terminal design primitives with zero external dependencies.

Part of the [Krewire](https://krewire.com) digital SDLC ecosystem.

## Features & Modules

Krewire TUI is organized as a multi-module monorepo with 2 independent Go modules:

- **Core Framework (`github.com/krewire/tui/core`)**: Structured application model (`App`), `Command`, `Flag`, subcommands, fuzzy suggestions, and progressive disclosure help formatting.
- **Terminal Primitives (`github.com/krewire/tui/term`)**: Pure-Go ANSI colors, styling (bold, dim, underline), terminal capability detection (`NO_COLOR`, `TERM=dumb`), stream routing, and escape sequence sanitization.

## Installation

```bash
go get github.com/krewire/tui/core
go get github.com/krewire/tui/term
```

## Architecture & Invariants

- **Pure Go & CGO-Free**: 100% standard library-first design compiling with `CGO_ENABLED=0`.
- **Zero Third-Party Dependencies**: Zero external dependencies outside the Go standard library and foundational `github.com/krewire/krewire/errs`.
- **Strict Downward Dependencies**: `term` is a leaf package and never imports `core`.
- **ANSI Escape Sanitization**: Output streams sanitize malicious CSI/OSC/C1 escapes to prevent terminal injection.

## Quality Gates & CI

Quality gates are enforced locally and in GitHub Actions CI:

```bash
# Code formatting
gofmt -l .

# Static analysis across modules
go vet ./core/... ./term/...

# Unit and race tests
go test -v -race ./core/... ./term/...

# Pure-Go compilation
CGO_ENABLED=0 go build ./core/... ./term/...
```

### CI Required Checks

- `gofmt`: Ensures all `.go` files follow standard Go formatting.
- `module-quality (core)`: Enforces vet, race tests, and pure-Go build for `core`.
- `module-quality (term)`: Enforces vet, race tests, and pure-Go build for `term`.
- `terminal-conformance`: Verifies terminal behavior, escape sanitization, and ANSI safety.
- `ci-ok`: Required aggregating status check. Any failure in constituent jobs blocks merge.

## Versioning & Release Policy

Krewire TUI is currently in pre-1.0 alpha (`v0.x`) phase as a first-party component of Krewire CLI (`krew`). Semver stability (`v1.0.0`) and tagged GitHub releases will be established once the developer driver API contract freezes. Until then, consumers consume via local workspace (`go.work`) or git commit SHAs.

## License

[MIT License](LICENSE) © 2026 Krewire Contributors
