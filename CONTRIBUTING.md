# Contributing to Krewire TUI

Guidelines for contributing to **Krewire TUI** (`github.com/krewire/tui`).

> All code, comments, docstrings, commit messages, and documentation are strictly in **English**.

## Invariants

1. **Standard Library First**: Zero external dependencies outside the standard library and `github.com/krewire/krewire/errs`.
2. **Pure Go**: Always compile with `CGO_ENABLED=0`.
3. **Module Structure**: Code lives under `core/` or `term/` modules.
4. **Strict Downward Flow**: `term` never imports `core`.

## Quality Gates

Before opening a pull request, ensure all 4 quality gates pass:

```bash
# 1. Formatting (must return 0 lines)
gofmt -l .

# 2. Static analysis
go vet ./core/... ./term/...

# 3. Unit and race tests
go test -v -race ./core/... ./term/...

# 4. Pure-Go compilation
CGO_ENABLED=0 go build ./core/... ./term/...
```

## CI Pipeline & Required Checks

Pull requests and pushes to `main` execute `.github/workflows/ci.yaml`:

- `gofmt`: Formatting check (`gofmt -l .`).
- `module-quality (core)`: Vet, race test, and build for `core`.
- `module-quality (term)`: Vet, race test, and build for `term`.
- `terminal-conformance`: ANSI escape safety and terminal capabilities test suite.
- `ci-ok`: Aggregate check requiring all jobs to succeed.

## Versioning Policy

TUI follows semantic versioning aligned with the Krewire toolchain. Pre-1.0 development is coordinated across the workspace via `go.work`. Tagged releases will follow API stability ratification.
