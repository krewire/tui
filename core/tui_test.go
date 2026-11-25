package tui

import (
	"flag"
	"io"
	"strings"
	"testing"

	"github.com/krewire/krewire/errs"
)

func successHandler(*flag.FlagSet) errs.ExitCode {
	return errs.ExitCodeSuccess
}

func testApp() *App {
	app := NewApp("demo", "0.1.0").
		Command(NewCommand("ping", "ping the app", nil, successHandler))
	app.stderr = io.Discard
	return app
}

func TestRunDispatchesToKnownCommand(t *testing.T) {
	app := testApp()
	if got := app.Run([]string{"ping"}); got != errs.ExitCodeSuccess {
		t.Errorf("Run() = %v, want %v", got, errs.ExitCodeSuccess)
	}
}

func TestRunEmptyIsUsage(t *testing.T) {
	app := testApp()
	if got := app.Run(nil); got != errs.ExitCodeUsage {
		t.Errorf("Run() = %v, want %v", got, errs.ExitCodeUsage)
	}
}

func TestRunUnknownCommandIsUsage(t *testing.T) {
	app := testApp()
	if got := app.Run([]string{"nope"}); got != errs.ExitCodeUsage {
		t.Errorf("Run() = %v, want %v", got, errs.ExitCodeUsage)
	}
}

func TestRunParseErrorIsUsage(t *testing.T) {
	app := NewApp("demo", "0.1.0").
		Command(NewCommand("ping", "ping the app", func(fs *flag.FlagSet) {
			fs.String("name", "", "name of the target")
		}, successHandler))
	app.stderr = io.Discard
	if got := app.Run([]string{"ping", "--nope"}); got != errs.ExitCodeUsage {
		t.Errorf("Run() = %v, want %v", got, errs.ExitCodeUsage)
	}
}

func TestRunHelpGeneral(t *testing.T) {
	app := testApp()
	app.stderr = io.Discard
	if got := app.Run([]string{"help"}); got != errs.ExitCodeSuccess {
		t.Errorf("Run(help) = %v, want %v", got, errs.ExitCodeSuccess)
	}
	if got := app.Run([]string{"--help"}); got != errs.ExitCodeSuccess {
		t.Errorf("Run(--help) = %v, want %v", got, errs.ExitCodeSuccess)
	}
	if got := app.Run([]string{"-h"}); got != errs.ExitCodeSuccess {
		t.Errorf("Run(-h) = %v, want %v", got, errs.ExitCodeSuccess)
	}
}

func TestRunHelpCommand(t *testing.T) {
	app := testApp()
	app.stderr = io.Discard
	if got := app.Run([]string{"help", "ping"}); got != errs.ExitCodeSuccess {
		t.Errorf("Run(help ping) = %v, want %v", got, errs.ExitCodeSuccess)
	}
	if got := app.Run([]string{"help", "unknown"}); got != errs.ExitCodeUsage {
		t.Errorf("Run(help unknown) = %v, want %v", got, errs.ExitCodeUsage)
	}
}

func TestRunCommandHelpAliases(t *testing.T) {
	app := testApp()
	app.stderr = io.Discard
	aliases := [][]string{
		{"ping", "help"},
		{"ping", "--help"},
		{"ping", "-h"},
		{"ping", "-help"},
	}
	for _, args := range aliases {
		if got := app.Run(args); got != errs.ExitCodeSuccess {
			t.Errorf("Run(%v) = %v, want %v", args, got, errs.ExitCodeSuccess)
		}
	}
}

func TestHiddenCommandCanRunAndIsHiddenFromGeneralHelp(t *testing.T) {
	called := false
	hiddenHandler := func(*flag.FlagSet) errs.ExitCode {
		called = true
		return errs.ExitCodeSuccess
	}
	app := NewApp("demo", "0.1.0").
		Command(NewCommand("secret", "secret easter egg", nil, hiddenHandler).WithHidden(true)).
		Command(NewCommand("ping", "ping command", nil, successHandler))
	app.stderr = io.Discard

	// Should still execute directly
	if got := app.Run([]string{"secret"}); got != errs.ExitCodeSuccess || !called {
		t.Errorf("Run(secret) = %v, called = %v", got, called)
	}

	// Should not be in general help listing
	groups := app.groupedCommands(false)
	for _, g := range groups {
		for _, cmd := range g.cmds {
			if cmd.Name == "secret" {
				t.Errorf("hidden command %q found in groupedCommands", cmd.Name)
			}
		}
	}
}

func TestCommandDispatch_TableDriven(t *testing.T) {
	var capturedName string
	var capturedVerbose bool
	var capturedCount int
	var capturedArgs []string

	app := NewApp("testapp", "1.0.0").
		Command(NewCommand("greet", "greet someone", func(fs *flag.FlagSet) {
			fs.StringVar(&capturedName, "name", "world", "person to greet")
			fs.BoolVar(&capturedVerbose, "verbose", false, "verbose output")
			fs.IntVar(&capturedCount, "count", 1, "repetition count")
		}, func(fs *flag.FlagSet) errs.ExitCode {
			capturedArgs = fs.Args()
			if capturedName == "error" {
				return errs.ExitCodeFailure
			}
			if capturedName == "missing" {
				return errs.ExitCodeUsage
			}
			return errs.ExitCodeSuccess
		})).
		WithHint("start", "testapp greet")
	app.stderr = io.Discard

	cases := []struct {
		name     string
		tokens   []string
		wantExit errs.ExitCode
		check    func() bool
	}{
		{
			name:     "default flags",
			tokens:   []string{"greet"},
			wantExit: errs.ExitCodeSuccess,
			check: func() bool {
				return capturedName == "world" && !capturedVerbose && capturedCount == 1
			},
		},
		{
			name:     "override flags",
			tokens:   []string{"greet", "--name", "alice", "--verbose", "--count", "3"},
			wantExit: errs.ExitCodeSuccess,
			check: func() bool {
				return capturedName == "alice" && capturedVerbose && capturedCount == 3
			},
		},
		{
			name:     "flag reordering with positional args",
			tokens:   []string{"greet", "extra1", "--name", "bob", "extra2"},
			wantExit: errs.ExitCodeSuccess,
			check: func() bool {
				return capturedName == "bob" && len(capturedArgs) == 2 && capturedArgs[0] == "extra1" && capturedArgs[1] == "extra2"
			},
		},
		{
			name:     "propagates failure exit code",
			tokens:   []string{"greet", "--name", "error"},
			wantExit: errs.ExitCodeFailure,
			check:    func() bool { return true },
		},
		{
			name:     "propagates usage exit code",
			tokens:   []string{"greet", "--name", "missing"},
			wantExit: errs.ExitCodeUsage,
			check:    func() bool { return true },
		},
		{
			name:     "invalid flag returns usage",
			tokens:   []string{"greet", "--unrecognized"},
			wantExit: errs.ExitCodeUsage,
			check:    func() bool { return true },
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			capturedName = ""
			capturedVerbose = false
			capturedCount = 0
			capturedArgs = nil

			got := app.Run(tc.tokens)
			if got != tc.wantExit {
				t.Errorf("Run(%v) exit = %v, want %v", tc.tokens, got, tc.wantExit)
			}
			if !tc.check() {
				t.Errorf("check failed for %v", tc.tokens)
			}
		})
	}
}

func TestApp_StreamIsolation(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}

	app := NewApp("streamapp", "1.0.0").
		Command(NewCommand("export", "export data", nil, func(fs *flag.FlagSet) errs.ExitCode {
			stdout.WriteString(`{"status":"ok"}`)
			return errs.ExitCodeSuccess
		}))
	app.stderr = stderr

	// Run command that produces stdout
	code := app.Run([]string{"export"})
	if code != errs.ExitCodeSuccess {
		t.Fatalf("export failed: %v", code)
	}
	if stdout.String() != `{"status":"ok"}` {
		t.Errorf("stdout corrupted: %q", stdout.String())
	}

	// Trigger help/usage: must write to stderr, not polluting stdout
	stdout.Reset()
	stderr.Reset()
	helpCode := app.Run([]string{"--help"})
	if helpCode != errs.ExitCodeSuccess {
		t.Fatalf("help failed: %v", helpCode)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected clean stdout on help, got: %q", stdout.String())
	}
	if stderr.Len() == 0 {
		t.Errorf("expected stderr to contain help content, got empty")
	}
}
