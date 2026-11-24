# Krewire TUI — Architecture & Subsystems Design

This document details the architectural boundaries, design invariants, and package structure of `github.com/krewire/tui`.

## Subsystem Architecture

The TUI repository is structured into two focused packages:

```text
                  github.com/krewire/krew (CLI Toolchain)
                                    │
                                    ▼
                         github.com/krewire/tui
                  - App & Command Group Routing
                  - Flag Parsing & Options
                  - Help Output & Suggestions
                                    │
                                    ▼
                       github.com/krewire/tui/term (Leaf)
                  - ANSI Escape Sequences & Color Painting
                  - Terminal Dimensions & Detection
                  - Word Wrapping & Text Formatting
                                    │
                                    ▼
                         Go Standard Library (os, io, fmt)
```

## Package Responsibilities

### 1. `github.com/krewire/tui`
- **Application Builder (`App`)**: Encapsulates application metadata (name, version, description), command routing, and error exit codes.
- **Command Dispatcher (`Command`)**: Manages command hierarchies, flag sets, aliases, examples, and grouping (`WithGroup`, `WithCore`).
- **Help Formatter**: Generates human-friendly, ANSI-formatted help listings with command categories and did-you-mean typo suggestions.

### 2. `github.com/krewire/tui/term`
- **Terminal Primitives (`Terminal`)**: ANSI 16-color, 256-color, and TrueColor escape codes.
- **Style Modifiers**: Bold, Dim, Italic, Underline, and Invert.
- **Word Wrapping & Layout**: Bounds checking, word wrapping to terminal width, and line truncation.
- **Terminal Injection Protection (`Sanitize`)**: Untrusted caller text passed to `Paint` is sanitized by default against CSI, OSC, BEL, and C0/C1 control characters. Callers requiring verbatim escape sequences must explicitly invoke `RawPaint`.
- **Zero-Dependency Leaf Invariant**: Must never import `github.com/krewire/tui` or any package outside the Go standard library.

## Architectural Invariants

1. **Leaf Independence**: `term` must remain an independent, zero-dependency leaf package. This invariant is enforced by `architecture_test.go`.
2. **Pure Go (`CGO_ENABLED=0`)**: The framework must compile cleanly without CGO on all supported OS targets.
3. **No External Dependencies**: `go.mod` requires zero external third-party dependencies.
