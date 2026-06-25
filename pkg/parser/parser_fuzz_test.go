// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/lexer"
)

func FuzzParser(f *testing.F) {
	testcases := []string{
		"echo hello",
		"if [[ -f foo ]]; then echo bar; fi",
		"ls -la | grep foo",
		"${var:-default}",
		"$(command)",
		"for i in 1 2 3; do echo $i; done",
		"case $foo in bar) echo match ;; esac",
		"func() { echo body; }",
		"time command",
		"coproc command",
		// Command-position assignments exercise reshapeCommandAssignment,
		// including a multi-part RHS and a truncated form whose parts may
		// be incomplete.
		"foo && cd=$bar",
		"a | b=$c/$d",
		"x && y=$(",
	}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, input string) {
		l := lexer.New(input)
		p := New(l)
		p.ParseProgram()
	})
}
