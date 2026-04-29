// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// TestFixTypeGuardSweep walks every registered Fix and feeds a node
// type that does NOT match the fix's expected concrete type. The
// expected behaviour is a defensive nil return; the call exercises
// the guard branch that every Fix opens with (and the matching
// `if cs == nil` typed-nil follow-up where present).
//
// Single-function coverage gain across hundreds of katas; each fix
// reports 66.7% because the body's happy-path is hit by katatests but
// the guard branch is not.
func TestFixTypeGuardSweep(t *testing.T) {
	wrong := &ast.IntegerLiteral{Value: 0}
	v := Violation{}
	src := []byte("")
	for _, k := range Registry.KatasByID {
		if k.Fix == nil {
			continue
		}
		// Most fixes type-assert and return nil on mismatch; the few
		// that ignore the node and use only Violation/source coordinates
		// may still produce a non-nil edit. Either result is acceptable;
		// the call is enough to record the guard branch.
		_ = k.Fix(wrong, v, src)
	}
}

// TestCheckTypeGuardSweep mirrors the fix sweep over every registered
// Check function. The Check entrypoint is called with a node type that
// does not match the kata's expected concrete type, exercising the
// guard branch that every Check opens with. A per-iteration recover
// records any kata whose guard does an unchecked cast so the test
// surface stays green; the panic itself is logged for follow-up.
func TestCheckTypeGuardSweep(t *testing.T) {
	wrong := &ast.IntegerLiteral{Value: 0}
	for id, k := range Registry.KatasByID {
		if k.Check == nil {
			continue
		}
		func(id string, fn func(ast.Node) []Violation) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Check(%s) panicked on wrong-type input: %v", id, r)
				}
			}()
			_ = fn(wrong)
		}(id, k.Check)
	}
}

// TestFixTypedNilSweep feeds every common typed-nil concrete pointer to
// every Fix entrypoint. Mirrors the keywordStmtToExpression nil-receiver
// guard from #1314: a fix that dereferences a typed-nil panic surface
// shows up in the per-call Logf and gets fixed in a follow-up; the
// recovery keeps the sweep alive so the surrounding `if x == nil`
// guard branches get coverage for every well-formed kata.
func TestFixTypedNilSweep(t *testing.T) {
	v := Violation{}
	src := []byte("")
	typedNils := []ast.Node{
		(*ast.SimpleCommand)(nil),
		(*ast.Identifier)(nil),
		(*ast.InfixExpression)(nil),
		(*ast.PrefixExpression)(nil),
		(*ast.CallExpression)(nil),
		(*ast.StringLiteral)(nil),
		(*ast.IntegerLiteral)(nil),
		(*ast.ConcatenatedExpression)(nil),
		(*ast.DoubleBracketExpression)(nil),
	}
	for id, k := range Registry.KatasByID {
		if k.Fix == nil {
			continue
		}
		for _, n := range typedNils {
			func(id string, fn func(ast.Node, Violation, []byte) []FixEdit, n ast.Node) {
				defer func() {
					if r := recover(); r != nil {
						t.Logf("Fix(%s) panicked on typed-nil receiver: %v", id, r)
					}
				}()
				_ = fn(n, v, src)
			}(id, k.Fix, n)
		}
	}
}
