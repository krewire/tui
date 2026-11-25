# `core` Module Reference

- **Module Path**: `github.com/krewire/tui/core` (`core/go.mod`)
- **Package Name**: `package tui`

## Purpose & Domain Responsibility

`core` provides the command-line application model and TUI orchestration layer for Krewire CLI tools. It composes Go's standard `flag.FlagSet` and `log/slog` with `github.com/krewire/tui/term` and `github.com/krewire/krewire/errs` to deliver subcommand routing, flexible flag reordering, progressive-disclosure help, typed `<APP>_` environment configuration, machine-readable JSON output, and typo suggestions.

## Architectural Position

- **Layer**: Layer 2 (Core TUI / CLI Framework)
- **Pure-Go Invariant**: Compiles cleanly with `CGO_ENABLED=0`
- **Dependencies**:
  - Third-party: None
  - Ecosystem: `github.com/krewire/tui/term`, `github.com/krewire/krewire/errs`

## Core Concepts & APIs

### Application & Command Routing (`tui.go`)

- **`ExitCode`**: Type alias for `errs.ExitCode` with constants `ExitCodeSuccess`, `ExitCodeFailure`, and `ExitCodeUsage`.
- **`Command`**: `{ Name, About, Group, Example, Hidden, Core bool, Register func(*flag.FlagSet), Run func(*flag.FlagSet) errs.ExitCode }`.
  - Constructed via **`NewCommand(name, about, register, run)`** and configured fluently with **`WithGroup(group)`**, **`WithExample(example)`**, **`WithHidden(hidden)`**, and **`WithCore()`**.
- **`App`**: Constructed via **`NewApp(name, version string) *App`** and configured with **`WithDescription(desc)`**, **`WithHint(name, suggestion)`**, and **`Command(cmd)`**.
  - **`Name() string`**: Returns the application name.
  - **`Run(tokens []string) errs.ExitCode`**: Handles empty args (prints help and returns `ExitCodeUsage`), global `--help`/`-h`/`--all`/`-a`, `help [command]`, `<command> help|--help|-h`, unknown command diagnostics, and flag parsing.
  - **Flag Reordering**: Before calling `fs.Parse`, `App.Run` reorders flags placed after positional arguments (e.g., `mycli build ./dir --minify`) while respecting `--` end-of-options markers.

### Progressive-Disclosure Help (`help.go`)

When any registered command has `Core: true`, general help (`help` / `--help`) displays a `CORE WORKFLOW` section with full descriptions and examples, followed by `ADVANCED CAPABILITIES` listing remaining commands compactly by group (sorted deterministically: `project`, `inspect`, `build`, `develop`, `ship`, then custom/other). Passing `help --all` or `-a` expands every command in full.

### Typed Environment Configuration (`config.go`)

- **`Env`**: Constructed with **`NewEnv(appName string) *Env`**, normalizing `appName` to `SCREAMING_SNAKE_CASE_` (e.g., `"my-cli"` → `"MY_CLI_"`).
- **`Prefix() string`**: Returns the resolved prefix.
- **`GetString(key, def string) string`**, **`GetBool(key string, def bool) (bool, error)`**, and **`GetInt(key string, def int) (int, error)`**: Read `<PREFIX><KEY>` from `os.LookupEnv`, returning `def` when unset/empty and an error on invalid boolean/integer syntax.

### Structured JSON Output (`output.go`)

- **`RenderJSON(v any) ([]byte, error)`** and **`WriteJSON(w io.Writer, v any) error`**: Emit deterministic 2-space-indented JSON to a data stream (`os.Stdout`), keeping machine output separate from stderr diagnostics.
- **`Result`**: `{ OK bool, Data any, Error string }` with constructors **`ResultOK(data any)`** and **`ResultErr(msg string)`**.

### Typo Suggestions (`suggest.go`)

- **`Suggest(name string, candidates []string) string`**: Returns the first candidate where `name` is a case-insensitive prefix, or otherwise the closest candidate within Levenshtein edit distance `<= 2` (or `""` if none matches). `App.Run` checks `WithHint` overrides first before falling back to `Suggest`.

## Usage Example

```go
package main

import (
	"flag"
	"os"

	tui "github.com/krewire/tui/core"
)

func main() {
	env := tui.NewEnv("mycli")

	app := tui.NewApp("mycli", "1.0.0").
		WithDescription("Example CLI built on tui/core").
		WithHint("start", "mycli dev  (starts the dev server)").
		Command(tui.NewCommand("dev", "Start dev server",
			func(fs *flag.FlagSet) {
				fs.String("addr", env.GetString("ADDR", ":8080"), "Listen address")
			},
			func(fs *flag.FlagSet) tui.ExitCode {
				addr := fs.Lookup("addr").Value.String()
				_ = tui.WriteJSON(os.Stdout, tui.ResultOK(map[string]string{"addr": addr}))
				return tui.ExitCodeSuccess
			},
		).WithGroup("develop").WithCore().WithExample("mycli dev --addr :9090"))

	os.Exit(int(app.Run(os.Args[1:])))
}
```

## Verification

```bash
cd /home/krewire/workspace/krewire/tui/core
gofmt -l .
go vet ./...
CGO_ENABLED=0 go test -v ./...
CGO_ENABLED=0 go build ./...
```
