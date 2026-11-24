# Contributing to Krewire TUI

Guidelines for contributing to **Krewire TUI** (`github.com/krewire/tui`).

> All code, comments, docstrings, commit messages, and documentation are strictly in **English**.

## Invariants

1. **Standard Library Only**: Never introduce third-party dependencies.
2. **Pure Go**: Always compile with `CGO_ENABLED=0`.
3. **Module Structure**: Code lives under `core/` or `term/` modules.

## Quality Gates

```bash
gofmt -l .
go vet ./...
go test ./...
CGO_ENABLED=0 go build ./...
```
