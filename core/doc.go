// Package tui provides a lightweight, pure-Go Terminal User Interface (TUI)
// and command-line application model with zero external dependencies.
//
// # Architectural Role & Boundaries
//
// The tui repository sits at Layer 2 in the Krewire ecosystem pyramid. It is
// designed as a high-cohesion, minimalist foundation for developer tooling and
// interactive console applications.
//
// Responsibilities supported by the root package:
//   - Application lifecycle management (App, Run, Versioning, Slog integration).
//   - Command tree composition (Command, sub-commands, command groups, aliases).
//   - Progressive disclosure help (Core workflows vs. detailed command listing).
//   - Output formatting & streams (Data on stdout vs. diagnostics on stderr).
//   - Fuzzy/Levenshtein command suggestions for mistyped commands.
//
// # Package Structure & Modular Growth Boundaries
//
//   - github.com/krewire/tui (Root):
//     The primary public entrypoint and unified application model.
//
//   - github.com/krewire/tui/term:
//     Terminal capability detection, ANSI styling, and raw terminal primitives.
//     Strictly independent; never imports the root tui package.
//
// Additional subpackages (e.g., layout, prompt, table, widget) should only be
// introduced when a subsystem has a distinct API contract, stateful behavior,
// and can be tested independently. Avoid creating generic 'utils', 'common',
// or 'internal' packages without concrete ownership needs.
package tui
