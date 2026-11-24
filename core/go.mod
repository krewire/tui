module github.com/krewire/tui/core

go 1.27.3

require (
	github.com/krewire/krewire/errs v0.0.0
	github.com/krewire/tui/term v0.0.0
)

replace (
	github.com/krewire/krewire/errs => ../../krewire/errs
	github.com/krewire/tui/term => ../term
)
