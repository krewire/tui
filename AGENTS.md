# AGENTS.md — Krewire TUI

Operational guide for autonomous AI agents working inside the `krewire/tui` repository.

## Architecture

- Multi-module monorepo comprising `core/` and `term/`.
- No `go.mod` at repository root; local coordination via `go.work`.
- Strictly pure-Go standard library only (`CGO_ENABLED=0`).
- All code, comments, docstrings, and commit messages must be strictly in English.

## Quality Gates

Before finishing work:
```bash
gofmt -l .
go vet ./...
go test ./...
CGO_ENABLED=0 go build ./...
```
