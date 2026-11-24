package tui

import (
	"bytes"
	"flag"
	"strings"
	"testing"

	"github.com/krewire/krewire/errs"
	"github.com/krewire/tui/term"
)

func disclosureApp(buf *bytes.Buffer) *App {
	app := NewApp("demo", "0.1.0").
		Command(NewCommand("new", "create a project", nil, successHandler).WithGroup("project").WithCore().WithExample("demo new app")).
		Command(NewCommand("dev", "live server", nil, successHandler).WithGroup("develop").WithCore()).
		Command(NewCommand("init", "set up in place", nil, successHandler).WithGroup("project").WithExample("demo init --site")).
		Command(NewCommand("tier", "inspect tiers", nil, successHandler).WithGroup("inspect")).
		WithHint("start", "demo dev")
	app.stderr = buf
	app.terminal = &term.Terminal{}
	return app
}

func TestHelpProgressiveDisclosure(t *testing.T) {
	var buf bytes.Buffer
	app := disclosureApp(&buf)
	if code := app.Run([]string{"help"}); code != errs.ExitCodeSuccess {
		t.Fatalf("code %v", code)
	}
	out := buf.String()
	for _, want := range []string{"CORE WORKFLOW", "new   create a project", "→ demo new app", "ADVANCED CAPABILITIES", "PROJECT  init", "INSPECT  tier", "demo help --all"} {
		if !strings.Contains(out, want) {
			t.Errorf("general help missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "set up in place") || strings.Contains(out, "demo init --site") {
		t.Errorf("advanced commands should be listed by name only:\n%s", out)
	}
	if strings.Index(out, "CORE WORKFLOW") > strings.Index(out, "ADVANCED CAPABILITIES") {
		t.Error("core section must come first")
	}
}

func TestHelpAllListsEverything(t *testing.T) {
	for _, args := range [][]string{{"help", "--all"}, {"help", "-a"}, {"--help", "--all"}, {"--all"}} {
		var buf bytes.Buffer
		app := disclosureApp(&buf)
		if code := app.Run(args); code != errs.ExitCodeSuccess {
			t.Fatalf("%v: code %v", args, code)
		}
		out := buf.String()
		for _, want := range []string{"CORE WORKFLOW", "set up in place", "→ demo init --site", "inspect tiers"} {
			if !strings.Contains(out, want) {
				t.Errorf("%v missing %q:\n%s", args, want, out)
			}
		}
		if strings.Contains(out, "help --all to describe") {
			t.Errorf("%v: full reference should not point to --all", args)
		}
	}
}

func TestHelpWithoutCoreKeepsFlatGroups(t *testing.T) {
	var buf bytes.Buffer
	app := NewApp("demo", "0.1.0").Command(NewCommand("ping", "ping the app", nil, successHandler))
	app.stderr, app.terminal = &buf, &term.Terminal{}
	app.Run([]string{"help"})
	if out := buf.String(); !strings.Contains(out, "COMMANDS") || !strings.Contains(out, "ping  ping the app") || strings.Contains(out, "CORE WORKFLOW") {
		t.Errorf("help:\n%s", out)
	}
}

func TestUnknownCommandHint(t *testing.T) {
	var buf bytes.Buffer
	app := disclosureApp(&buf)
	if code := app.Run([]string{"start"}); code != errs.ExitCodeUsage {
		t.Fatalf("code %v", code)
	}
	if out := buf.String(); !strings.Contains(out, `unknown command "start"`) || !strings.Contains(out, "did you mean demo dev") {
		t.Errorf("output:\n%s", out)
	}
}

func TestSuggest(t *testing.T) {
	cands := []string{"build", "bundle", "serve"}
	for in, want := range map[string]string{"buidl": "build", "bun": "bundle", "SERVE": "serve", "zzzzzz": ""} {
		if got := Suggest(in, cands); got != want {
			t.Errorf("Suggest(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestHelpRenderer_SanitizesMaliciousMetadata(t *testing.T) {
	var buf bytes.Buffer

	// Malicious command metadata attempting to inject OSC 8 hyperlinks,
	// OSC 52 clipboard commands, CSI cursor line-wipe, and carriage returns
	maliciousApp := NewApp("app\x1b]8;;https://evil.test\x07name\x1b]8;;\x07", "1.0.0").
		Command(NewCommand("build", "build project\x1b[2K\r[OK] fake status", func(fs *flag.FlagSet) {
			fs.String("config", "", "path to config\x1b]52;c;evil\x07 file\roverwrite")
		}, successHandler).WithExample("build --config\x07 evil.json"))
	maliciousApp.stderr = &buf
	maliciousApp.terminal = &term.Terminal{ColorSupported: false}

	// 1. General help
	code := maliciousApp.Run([]string{"help"})
	if code != errs.ExitCodeSuccess {
		t.Fatalf("unexpected exit code: %v", code)
	}

	out := buf.String()
	if strings.Contains(out, "\x1b]8;;") {
		t.Errorf("general help output contains unescaped OSC 8 hyperlink: %q", out)
	}
	if strings.Contains(out, "\x1b[2K") {
		t.Errorf("general help output contains unescaped CSI line wipe: %q", out)
	}
	if strings.Contains(out, "\x07") {
		t.Errorf("general help output contains unescaped BEL byte: %q", out)
	}
	if strings.Contains(out, "\r") {
		t.Errorf("general help output contains unescaped isolated carriage return: %q", out)
	}

	// 2. Command specific help
	buf.Reset()
	code = maliciousApp.Run([]string{"build", "--help"})
	if code != errs.ExitCodeSuccess {
		t.Fatalf("unexpected exit code: %v", code)
	}

	cmdOut := buf.String()
	if strings.Contains(cmdOut, "\x1b]52;") {
		t.Errorf("command help output contains unescaped OSC 52 clipboard sequence: %q", cmdOut)
	}
	if strings.Contains(cmdOut, "\x07") {
		t.Errorf("command help output contains unescaped BEL byte: %q", cmdOut)
	}
	if strings.Contains(cmdOut, "\r") {
		t.Errorf("command help output contains unescaped isolated carriage return: %q", cmdOut)
	}
}
