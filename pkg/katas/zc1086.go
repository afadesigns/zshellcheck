// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.FunctionDefinitionNode, Kata{
		ID:    "ZC1086",
		Title: "Prefer `func() { ... }` over `function func { ... }`",
		Description: "The `function` keyword is optional in Zsh and non-standard in POSIX sh. " +
			"Using `func() { ... }` is more portable and consistent.",
		Severity: SeverityStyle,
		Check:    checkZC1086,
		Fix:      fixZC1086,
	})
	RegisterKata(ast.FunctionLiteralNode, Kata{
		ID:    "ZC1086",
		Title: "Prefer `func() { ... }` over `function func { ... }`",
		Description: "The `function` keyword is optional in Zsh and non-standard in POSIX sh. " +
			"Using `func() { ... }` is more portable and consistent.",
		Severity: SeverityStyle,
		Check:    checkZC1086,
		Fix:      fixZC1086,
	})
}

// fixZC1086 rewrites `function name [()] { body }` to the portable
// `name() { body }` form. Deletes the `function ` prefix and, when
// the source doesn't already carry `()` after the name, inserts it.
func fixZC1086(node ast.Node, v Violation, source []byte) []FixEdit {
	var name string
	switch n := node.(type) {
	case *ast.FunctionLiteral:
		if n.TokenLiteral() != "function" || n.Name == nil {
			return nil
		}
		name = n.Name.Value
	case *ast.FunctionDefinition:
		if n.TokenLiteral() != "function" || n.Name == nil {
			return nil
		}
		name = n.Name.Value
	default:
		return nil
	}
	if name == "" {
		return nil
	}
	kwOffset := LineColToByteOffset(source, v.Line, v.Column)
	if kwOffset < 0 || kwOffset+len("function ") > len(source) {
		return nil
	}
	if string(source[kwOffset:kwOffset+len("function")]) != "function" {
		return nil
	}
	// Delete `function` + the whitespace run that follows, up to the
	// name start. Find the name's byte position by scanning forward
	// past the whitespace.
	i := kwOffset + len("function")
	for i < len(source) && (source[i] == ' ' || source[i] == '\t') {
		i++
	}
	if i+len(name) > len(source) || string(source[i:i+len(name)]) != name {
		return nil
	}
	edits := []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  i - kwOffset,
		Replace: "",
	}}
	// After the name, check whether `()` already follows. Skip any
	// whitespace between name and peek byte; Zsh allows
	// `function foo ()` with a space.
	after := i + len(name)
	j := after
	for j < len(source) && (source[j] == ' ' || source[j] == '\t') {
		j++
	}
	if j >= len(source) || source[j] != '(' {
		afterLine, afterCol := offsetLineColZC1086(source, after)
		if afterLine < 0 {
			return nil
		}
		edits = append(edits, FixEdit{
			Line:    afterLine,
			Column:  afterCol,
			Length:  0,
			Replace: "()",
		})
	}
	return edits
}

func offsetLineColZC1086(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1086(node ast.Node) []Violation {
	// Case 1: function my_func { ... } -> Parsed as FunctionLiteralNode
	if funcLit, ok := node.(*ast.FunctionLiteral); ok {
		if funcLit.TokenLiteral() == "function" {
			return []Violation{
				{
					KataID:  "ZC1086",
					Message: "Prefer `func() { ... }` over `function func { ... }` for portability and consistency.",
					Line:    funcLit.TokenLiteralNode().Line,
					Column:  funcLit.TokenLiteralNode().Column,
					Level:   SeverityStyle,
				},
			}
		}
	}

	// Case 2: my_func() { ... } -> Parsed as FunctionDefinitionNode
	if funcDef, ok := node.(*ast.FunctionDefinition); ok {
		if funcDef.TokenLiteral() == "function" {
			return []Violation{
				{
					KataID:  "ZC1086",
					Message: "Prefer `func() { ... }` over `function func { ... }` for portability and consistency.",
					Line:    funcDef.TokenLiteralNode().Line,
					Column:  funcDef.TokenLiteralNode().Column,
					Level:   SeverityStyle,
				},
			}
		}
	}

	return nil
}
