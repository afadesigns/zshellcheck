// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

// TestWalkerTypedNilNoPanic guards against a regression where a kata's
// hand-rolled AST walker descended into a typed-nil interface value (an
// interface holding a nil concrete pointer, which is not == nil) and
// dereferenced a field on the nil receiver, crashing the whole run.
//
// The ZC1069 and ZC1053 walkers both lacked the reflect-based typed-nil
// guard that ast.Walk and the ZC1044 walker already carried. The inputs
// below are reduced from real files in the integration corpora that
// triggered SIGSEGV: prezto modules/terminal/init.zsh,
// zsh-syntax-highlighting main-highlighter.zsh, and the canonical Zsh
// Completion/Redhat/Command/_rpm. A panic here fails the test.
func TestWalkerTypedNilNoPanic(t *testing.T) {
	cases := []struct {
		name string
		src  string
	}{
		{
			name: "if cond subshell-negation backslash-continued",
			src:  "if [[ x == y ]] \\\n  && ( ! [[ -n \"$A\" ]] )\nthen\n  :\nfi\n",
		},
		{
			name: "if cond subshell-negation single line",
			src:  "if [[ $TERM == foo ]] && ( ! [[ -n \"$STY\" ]] ); then :; fi\n",
		},
		{
			name: "nested subshell in if condition",
			src:  "if ( ( ! [[ -n x ]] ) ); then local y=1; fi\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := lexer.New(tc.src)
			p := parser.New(l)
			program := p.ParseProgram()
			// Drive the full registered-kata surface over every node,
			// exactly as the CLI does. A typed-nil descent panics here.
			ast.Walk(program, func(n ast.Node) bool {
				_ = Registry.Check(n, nil)
				return true
			})
		})
	}
}
