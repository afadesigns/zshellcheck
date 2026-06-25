// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

// firstAutoload parses src and returns the first `autoload` command node
// plus the source bytes, for exercising the ZC1076 fix helpers directly.
func firstAutoload(t *testing.T, src string) (*ast.SimpleCommand, []byte) {
	t.Helper()
	prog := parser.New(lexer.New(src)).ParseProgram()
	var cmd *ast.SimpleCommand
	ast.Walk(prog, func(n ast.Node) bool {
		if c, ok := n.(*ast.SimpleCommand); ok && cmd == nil &&
			c.Name != nil && c.Name.String() == "autoload" {
			cmd = c
		}
		return true
	})
	if cmd == nil {
		t.Fatalf("no autoload command in %q", src)
	}
	return cmd, []byte(src)
}

// TestZC1076MissingFlags covers every branch, including the
// both-present case the Check -> Fix pipeline never reaches (the
// detector fires only when a flag is absent).
func TestZC1076MissingFlags(t *testing.T) {
	cases := []struct {
		src  string
		want string
	}{
		{"autoload foo", " -Uz"},
		{"autoload -U foo", " -z"},
		{"autoload -z foo", " -U"},
		{"autoload -Uz foo", ""},
		{"autoload -RUz foo", ""},
	}
	for _, c := range cases {
		cmd, _ := firstAutoload(t, c.src)
		if got := zc1076MissingFlags(cmd); got != c.want {
			t.Errorf("%q: got %q, want %q", c.src, got, c.want)
		}
	}
}

// TestFixZC1076BothFlagsPresent guards the defensive no-edit path: even
// if the fix is invoked directly on an already-correct command, it emits
// nothing rather than a duplicate flag.
func TestFixZC1076BothFlagsPresent(t *testing.T) {
	cmd, src := firstAutoload(t, "autoload -Uz foo")
	v := Violation{KataID: "ZC1076", Line: 1, Column: 1}
	if edits := fixZC1076(cmd, v, src); edits != nil {
		t.Errorf("expected no edits for an already-flagged autoload, got %v", edits)
	}
}
