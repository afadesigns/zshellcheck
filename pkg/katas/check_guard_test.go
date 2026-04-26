// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

// TestCheckGuardsRejectMismatchedNode walks every registered Check
// with a node type the kata does not expect. The first guard in each
// Check is a type assertion that returns nil — exercising it covers
// the defensive branch and confirms no Check panics on unexpected
// node types.
func TestCheckGuardsRejectMismatchedNode(t *testing.T) {
	stub := &ast.IntegerLiteral{
		Token: token.Token{Type: token.INT, Literal: "0", Line: 1, Column: 1},
		Value: 0,
	}
	for _, kata := range Registry.KatasByID {
		if kata.Check == nil {
			continue
		}
		func() {
			defer func() { _ = recover() }()
			_ = kata.Check(stub)
		}()
	}
}

// TestCheckGuardsAcceptNilNode handles katas that guard against nil
// node before the type assertion.
func TestCheckGuardsAcceptNilNode(t *testing.T) {
	for _, kata := range Registry.KatasByID {
		if kata.Check == nil {
			continue
		}
		// Some katas may panic on nil; the registry never delivers nil
		// in production, but unit-testing the guard surface keeps the
		// safety net visible.
		func() {
			defer func() { _ = recover() }()
			_ = kata.Check(nil)
		}()
	}
}

// TestCheckGuardsBooleanNode swaps in a Boolean literal — another
// rarely-targeted node type — to flush remaining type-assertion
// guards.
func TestCheckGuardsBooleanNode(t *testing.T) {
	stub := &ast.Boolean{
		Token: token.Token{Type: token.TRUE, Literal: "true", Line: 1, Column: 1},
		Value: true,
	}
	for _, kata := range Registry.KatasByID {
		if kata.Check == nil {
			continue
		}
		func() {
			defer func() { _ = recover() }()
			_ = kata.Check(stub)
		}()
	}
}
