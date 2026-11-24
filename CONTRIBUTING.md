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


## Continuous Integration

Pull requests targeting main and pushes to main run the GitHub Actions workflow in .github/workflows/ci.yml. The workflow is read-only, does not use production credentials, and records the commit SHA and selected Go toolchain in the job summary.

The blocking check names are:

- quality / Go quality gates
- terminal-conformance / terminal behavior

Both checks must pass before a change is considered verified. The Go quality gate runs gofmt validation, go vet, package tests, the race detector, and a CGO_ENABLED=0 build. The terminal conformance check runs the term package tests with a bounded timeout and non-interactive terminal settings. A skipped, cancelled, or missing check is not evidence of success.

For local verification, run:

    test -z "$(gofmt -l .)"
    go vet ./...
    go test -count=1 -timeout=3m ./...
    go test -race -count=1 -timeout=5m ./...
    CGO_ENABLED=0 go build ./...
    go test -count=1 -timeout=1m ./term

## Versioning and Releases

The compatibility and release policy is maintained in docs/release-policy.md. TUI remains pre-1.0 until maintainers explicitly approve a stable compatibility contract. Do not publish a stable v1.0.0 tag until the criteria in that policy have been reviewed and met.
