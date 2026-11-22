# Krewire TUI (`github.com/krewire/tui`)

Lightweight, pure-Go Terminal User Interface (TUI) framework and terminal design primitives with zero external dependencies.

Part of the [Krewire](https://krewire.com) ecosystem.

## Features

- **TUI Application Model (`tui`)**: Structured `App`, `Command`, `Flag`, subcommands, core workflow grouping, and progressive disclosure help.
- **Terminal Primitives (`term`)**: ANSI colors, text styles (bold, dim, underline), terminal capability detection, and output stream routing.
- **Zero NPM / Pure Go**: Standard library-first design compiling with `CGO_ENABLED=0`.

## Installation

```bash
go get github.com/krewire/tui
```

## Architecture & Boundaries

`tui` adheres to the Krewire 3S (Secure, Sustainable, Scalable) filter and Layer 2 modular design:

- **Root (`github.com/krewire/tui`)**: Unified application model, lifecycle management, sub-command router, flag parsing, progressive disclosure help formatting, and fuzzy suggestions.
- **Terminal Adapter (`github.com/krewire/tui/term`)**: Pure-Go ANSI color, styling, and terminal capability detection (`NO_COLOR`, `TERM=dumb`). Operates strictly as a leaf package and never imports root `tui`.
- **Modular Growth Guidelines**:
  - Keep the root core small and cohesive.
  - Subsystems with stateful APIs (such as interactive layout or widget engines) should only be split into dedicated packages when they have distinct independent lifecycles.
  - Avoid vague buckets like `utils` or `common`.

## Quick Start

```go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/krewire/tui"
	"github.com/krewire/tui/term"
)

func main() {
	app := tui.NewApp("mycli", "v1.0.0").
		WithDescription("Sample TUI tool").
		Command(tui.NewCommand("hello", "prints a greeting", nil, func(fs *flag.FlagSet) tui.ExitCode {
			tm := term.NewTerminal()
			fmt.Println(tm.Paint("Hello from Krewire TUI!", term.ColorGreen, []term.Style{term.StyleBold}))
			return tui.ExitCodeSuccess
		}).WithCore())

	os.Exit(int(app.Run(os.Args[1:])))
}
```

## License

MIT License. See [LICENSE](./LICENSE).
