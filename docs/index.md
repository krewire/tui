# Krewire TUI — Developer & Contributor Documentation

Welcome to the internal engineering documentation for **`krewire/tui`** (`github.com/krewire/tui`), the lightweight, pure-Go terminal user interface toolkit and ANSI styling engine for the Krewire digital SDLC ecosystem.

## Documentation Index

- [Architecture & Design](architecture.md): Two-package design (`tui` and `term`), leaf invariant, downward dependency rules, and terminal primitives.
- [Getting Started](getting-started.md): Installation, quickstart examples, command dispatching, flag registration, and formatting.
- [Contributing Guide](contributing.md): Quality gates, pure-Go invariant, testing standards, and pull request workflow.

## Overview & Core Principles

- **Zero External Dependencies**: Implemented strictly using the Go standard library (`os`, `fmt`, `strings`, `strconv`, `io`, `unicode`).
- **Pure Go (`CGO_ENABLED=0`)**: Completely CGO-free, enabling seamless cross-compilation across Linux, macOS, and Windows.
- **Two-Package Structure**:
  - `github.com/krewire/tui`: Top-level CLI command routing, flag binding, structured help generation, and typo suggestions.
  - `github.com/krewire/tui/term`: Leaf ANSI styling, color painting, screen dimension detection, and word wrapping.
- **Leaf Invariant**: `term` must never import `tui` or any package outside the Go standard library, enforced by `architecture_test.go`.
