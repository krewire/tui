# Krewire TUI (`github.com/krewire/tui`)

Lightweight, pure-Go Terminal User Interface (TUI) framework and terminal design primitives with zero external dependencies.

Part of the [Krewire](https://krewire.com) digital SDLC ecosystem.

## Features & Modules

Krewire TUI is organized as a multi-module monorepo with 2 independent Go modules:

- **Core Framework (`github.com/krewire/tui/core`)**: Structured application model (`App`), `Command`, `Flag`, subcommands, fuzzy suggestions, and progressive disclosure help formatting.
- **Terminal Primitives (`github.com/krewire/tui/term`)**: Pure-Go ANSI colors, styling (bold, dim, underline), terminal capability detection (`NO_COLOR`, `TERM=dumb`), and stream routing.

## Installation

```bash
go get github.com/krewire/tui/core
go get github.com/krewire/tui/term
```

## Architecture & Invariants

- **Pure Go & CGO-Free**: 100% standard library-first design compiling with `CGO_ENABLED=0`.
- **Zero Third-Party Dependencies**: Zero external dependencies outside the standard library.
- **Strict Downward Dependencies**: `term` is a leaf package and never imports `core`.

## Quality Gates

```bash
gofmt -l .
go vet ./...
go test ./...
CGO_ENABLED=0 go build ./...
```

## License

[MIT License](LICENSE) © 2026 Krewire Contributors
