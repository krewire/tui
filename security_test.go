package tui

import (
	"bytes"
	"flag"
	"strings"
	"testing"
	"unicode"

	"github.com/krewire/krewire/pkg/errs"
	"github.com/krewire/tui/term"
)

func TestHelpOutputSanitizesUntrustedMetadata(t *testing.T) {
	malicious := "\x1b[2Jvisible " +
		"\x1b]8;;https://evil.example\x07linked\x1b]8;;\x07 " +
		"\x1b]52;c;Ym9keQ==\x07 clipboard " +
		"\x07 BEL \r CR \u009b31m C1"

	var stderr bytes.Buffer
	app := NewApp("demo"+malicious, "1.0.0").WithDescription("description " + malicious)
	cmd := NewCommand("inspect", "about "+malicious, func(fs *flag.FlagSet) {
		fs.String("label", "", "flag usage "+malicious)
	}, func(*flag.FlagSet) errs.ExitCode {
		return errs.ExitCodeSuccess
	}).WithExample("demo inspect " + malicious)
	app.Command(cmd)
	app.stderr = &stderr
	app.terminal = &term.Terminal{ColorSupported: false}

	if got := app.Run([]string{"help", "inspect"}); got != errs.ExitCodeSuccess {
		t.Fatalf("Run(help inspect) = %v, want success", got)
	}

	out := stderr.String()
	for _, forbidden := range []string{"\x1b", "\x07", "\r", "\u009b", "evil.example", "Ym9keQ=="} {
		if strings.Contains(out, forbidden) {
			t.Errorf("help output contains forbidden terminal content %q: %q", forbidden, out)
		}
	}
	for _, want := range []string{"demo", "visible", "linked", "clipboard", "flag usage", "C1", "USAGE"} {
		if !strings.Contains(out, want) {
			t.Errorf("help output missing safe text %q: %s", want, out)
		}
	}
	for _, r := range out {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			t.Errorf("help output retained control rune U+%04X", r)
		}
	}
}

func TestRunParseErrorSanitizesUntrustedFlagValues(t *testing.T) {
	var stderr bytes.Buffer
	app := NewApp("demo", "1.0.0").Command(NewCommand("ping", "ping", func(fs *flag.FlagSet) {
		fs.Bool("enabled", false, "enable feature")
	}, func(*flag.FlagSet) errs.ExitCode {
		return errs.ExitCodeSuccess
	}))
	app.stderr = &stderr

	value := "bad\x1b[2Jvalue\x1b]52;c;Ym9keQ==\x07\r\u009b31m"
	if got := app.Run([]string{"ping", "--enabled=" + value}); got != errs.ExitCodeUsage {
		t.Fatalf("Run() = %v, want usage error", got)
	}
	out := stderr.String()
	for _, forbidden := range []string{"\x1b", "\x07", "\r", "\u009b", "Ym9keQ=="} {
		if strings.Contains(out, forbidden) {
			t.Errorf("parse diagnostic contains forbidden terminal content %q: %q", forbidden, out)
		}
	}
	if !strings.Contains(out, "badvalue") {
		t.Errorf("parse diagnostic lost safe text: %q", out)
	}
}
