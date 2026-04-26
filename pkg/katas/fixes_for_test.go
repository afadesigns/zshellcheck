// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func TestFixesForUnknownID(t *testing.T) {
	v := Violation{KataID: "ZC_NONEXISTENT", Line: 1, Column: 1}
	if got := Registry.FixesFor(nil, v, []byte{}); got != nil {
		t.Errorf("expected nil for unknown kata, got %v", got)
	}
}

func TestFixesForKataWithoutFix(t *testing.T) {
	kr := NewKatasRegistry()
	kr.RegisterKata(ast.IdentifierNode, Kata{
		ID:    "ZC_TEST_NOFIX",
		Title: "no fix",
		Check: func(ast.Node) []Violation { return nil },
	})
	v := Violation{KataID: "ZC_TEST_NOFIX", Line: 1, Column: 1}
	if got := kr.FixesFor(nil, v, []byte{}); got != nil {
		t.Errorf("expected nil for kata without Fix, got %v", got)
	}
}

func TestFixesForKataWithFix(t *testing.T) {
	called := false
	kr := NewKatasRegistry()
	kr.RegisterKata(ast.IdentifierNode, Kata{
		ID:    "ZC_TEST_WITHFIX",
		Title: "with fix",
		Check: func(ast.Node) []Violation { return nil },
		Fix: func(ast.Node, Violation, []byte) []FixEdit {
			called = true
			return []FixEdit{{Line: 1, Column: 1, Length: 0, Replace: "x"}}
		},
	})
	node := &ast.Identifier{Token: token.Token{Literal: "x"}, Value: "x"}
	v := Violation{KataID: "ZC_TEST_WITHFIX", Line: 1, Column: 1}
	got := kr.FixesFor(node, v, []byte("x"))
	if !called {
		t.Error("expected Fix to be invoked")
	}
	if len(got) != 1 || got[0].Replace != "x" {
		t.Errorf("unexpected fix output: %v", got)
	}
}
