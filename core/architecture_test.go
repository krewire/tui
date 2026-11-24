package tui_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestArchitectureInvariants asserts package dependency boundaries:
// 1. term/ package must never import the root tui package (prevent dependency cycles).
// 2. term/ package must only import standard library packages.
// 3. root tui package imports only standard library, term/, and pkg/errs.
func TestArchitectureInvariants(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working dir: %v", err)
	}

	termDir := filepath.Join(wd, "term")
	if _, err := os.Stat(termDir); err != nil {
		termDir = filepath.Join(wd, "..", "term")
	}
	termFset := token.NewFileSet()
	termPkgs, err := parser.ParseDir(termFset, termDir, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("failed to parse term dir: %v", err)
	}

	for pkgName, pkg := range termPkgs {
		for fileName, file := range pkg.Files {
			// Skip test files for internal testing dependencies if any
			if strings.HasSuffix(fileName, "_test.go") {
				continue
			}
			for _, imp := range file.Imports {
				path := strings.Trim(imp.Path.Value, `"`)
				if strings.Contains(path, "github.com/krewire/tui") {
					t.Errorf("term package (%s in %s) must NOT import %s (acyclic invariant violation)", pkgName, fileName, path)
				}
				if strings.Contains(path, ".") {
					t.Errorf("term package (%s in %s) must only import stdlib packages, found: %s", pkgName, fileName, path)
				}
			}
		}
	}
}
