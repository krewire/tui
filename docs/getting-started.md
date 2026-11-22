# Getting Started with Krewire TUI

This guide explains how to install and use `github.com/krewire/tui` in your Go applications and CLI tools.

## Installation

```bash
go get github.com/krewire/tui
```

## Quickstart Example

Create a simple CLI application with commands, flags, and ANSI terminal styling:

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
		WithDescription("Demo CLI powered by Krewire TUI").
		Command(tui.NewCommand("hello", "prints a greeting", registerHelloFlags, runHello).
			WithExample("mycli hello --name World"))

	os.Exit(int(app.Run(os.Args[1:])))
}

func registerHelloFlags(fs *flag.FlagSet) {
	fs.String("name", "World", "Name of person to greet")
}

func runHello(fs *flag.FlagSet) tui.ExitCode {
	name := "World"
	if f := fs.Lookup("name"); f != nil {
		name = f.Value.String()
	}

	tm := term.NewTerminal()
	cyan := tm.Paint(name, term.ColorCyan, []term.Style{term.StyleBold})
	fmt.Printf("Hello, %s!\n", cyan)

	return tui.ExitCodeSuccess
}
```

## Running Quality Gates

To verify formatting, static analysis, and tests:

```bash
cd tui
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
