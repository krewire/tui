# TUI Module Catalog

`tui` is organized as a two-module Go workspace (`tui/go.work`, with no root `go.mod`):

| Module | Import Path | Package Name | Focus |
|---|---|---|---|
| [`core`](./core.md) | `github.com/krewire/tui/core` | `tui` | CLI application builder (`App`), command dispatcher (`Command`), progressive-disclosure help, `<APP>_` environment config (`Env`), deterministic JSON output (`RenderJSON`, `Result`), and Levenshtein typo suggestions (`Suggest`). |
| [`term`](./term.md) | `github.com/krewire/tui/term` | `term` | Zero-dependency leaf package providing ANSI 16-color/style painting (`Paint`, `RawPaint`), terminal injection sanitization (`Sanitize`), TTY/environment capability detection (`NewTerminal`, `IsTerminal`), and stdout/stderr stream separation. |

## Dependency Boundary

```text
github.com/krewire/tui/core  (package tui)
   ├── imports github.com/krewire/tui/term
   └── imports github.com/krewire/krewire/errs
          │
          ▼
github.com/krewire/tui/term  (package term — Leaf)
   └── imports Go standard library only (enforced by core/architecture_test.go)
```
