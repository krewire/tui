# Contributing to Krewire TUI

Guidelines and engineering standards for contributors working on `github.com/krewire/tui`.

## Invariant Rules

1. **Pure Go**: All code must compile with `CGO_ENABLED=0`. Do not introduce CGO dependencies.
2. **Zero Dependencies**: Do not add third-party dependencies to `go.mod`. Use the Go standard library exclusively.
3. **Leaf Layering**: The `term` package must remain a leaf package. Never import `github.com/krewire/tui` or higher-level libraries from `term`.

## Testing Standards

- Write unit tests for all public methods and formatting helpers.
- Run `architecture_test.go` to ensure layering invariants are preserved.

## Verification Checklist

Before opening a pull request or committing changes:

```bash
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
