// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// TestZC1075HarvestHelpers covers the array-harvest helper guards,
// including the defensive type-assertion branches the corpus rarely hits.
func TestZC1075HarvestHelpers(t *testing.T) {
	if zc1075IsDeclName(&ast.StringLiteral{Value: "x"}) {
		t.Error("a non-identifier command name is not a declaration")
	}
	if zc1075IsDeclName(&ast.Identifier{Value: "echo"}) {
		t.Error("echo is not a declaration builtin")
	}
	if !zc1075IsDeclName(&ast.Identifier{Value: "local"}) {
		t.Error("local is a declaration builtin")
	}

	if _, ok := zc1075ArrayAssignName(&ast.Identifier{Value: "flags"}); ok {
		t.Error("a bare identifier is not an array assignment")
	}
	if _, ok := zc1075ArrayAssignName(&ast.ConcatenatedExpression{
		Parts: []ast.Expression{&ast.StringLiteral{Value: "("}},
	}); ok {
		t.Error("fewer than two parts is not an array assignment")
	}
	if _, ok := zc1075ArrayAssignName(&ast.ConcatenatedExpression{
		Parts: []ast.Expression{&ast.StringLiteral{Value: "x"}, &ast.ArrayLiteral{}},
	}); ok {
		t.Error("a non-identifier first part is not an array assignment")
	}
	name, ok := zc1075ArrayAssignName(&ast.ConcatenatedExpression{
		Parts: []ast.Expression{&ast.Identifier{Value: "arr"}, &ast.StringLiteral{Value: "="}, &ast.ArrayLiteral{}},
	})
	if !ok || name != "arr" {
		t.Errorf("array assignment should yield arr, got %q ok=%v", name, ok)
	}

	a := map[string]bool{}
	zc1075HarvestArrayDecl(&ast.DeclarationStatement{Flags: nil}, a)
	if len(a) != 0 {
		t.Error("a declaration with no flags harvests nothing")
	}

	if zc1075BareName("") != "" {
		t.Error("empty input yields empty name")
	}
}
