// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ArithmeticCommandNode, Kata{
		ID:          "ZC1073",
		Title:       "Unnecessary use of `$` in arithmetic expressions",
		Description: "Variables in `((...))` do not need `$` prefix. Use `(( var > 0 ))` instead of `(( $var > 0 ))`.",
		Severity:    SeverityStyle,
		Check:       checkZC1073,
		Fix:         fixZC1073,
	})
}

// fixZC1073 deletes the leading `$` from a variable used inside
// `(( … ))`. The violation coordinates already point at the `$`
// byte, so a single zero-replacement edit removes it. A second pass
// won't re-trigger because the identifier no longer carries `$`.
func fixZC1073(_ ast.Node, v Violation, source []byte) []FixEdit {
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off >= len(source) || source[off] != '$' {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  1,
		Replace: "",
	}}
}

func checkZC1073(node ast.Node) []Violation {
	cmd, ok := node.(*ast.ArithmeticCommand)
	if !ok {
		return nil
	}

	if cmd.Expression == nil {
		return nil
	}

	var violations []Violation

	ast.Walk(cmd.Expression, func(n ast.Node) bool {
		// Check for PrefixExpression with '$'
		if prefix, ok := n.(*ast.PrefixExpression); ok && prefix.Operator == "$" {
			if ident, ok := prefix.Right.(*ast.Identifier); ok {
				if isUserVariable(ident.Value) {
					violations = append(violations, Violation{
						KataID:  "ZC1073",
						Message: "Unnecessary use of `$` in arithmetic expressions. Use `(( var ))` instead of `(( $var ))`.",
						Line:    prefix.Token.Line,
						Column:  prefix.Token.Column,
						Level:   SeverityStyle,
					})
				}
			}
			return true
		}

		// Check for Identifier starting with '$' (if lexer emits VARIABLE)
		if ident, ok := n.(*ast.Identifier); ok {
			if len(ident.Value) > 1 && ident.Value[0] == '$' {
				varName := ident.Value[1:]
				if isUserVariable(varName) {
					violations = append(violations, Violation{
						KataID:  "ZC1073",
						Message: "Unnecessary use of `$` in arithmetic expressions. Use `(( var ))` instead of `(( $var ))`.",
						Line:    ident.Token.Line,
						Column:  ident.Token.Column,
						Level:   SeverityStyle,
					})
				}
			}
		}

		return true
	})

	return violations
}

func isUserVariable(name string) bool {
	if len(name) == 0 {
		return false
	}

	first := name[0]
	if !isAlpha(first) && first != '_' {
		return false
	}

	for i := 1; i < len(name); i++ {
		if !isAlphaNumeric(name[i]) && name[i] != '_' {
			return false
		}
	}

	return true
}

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isAlphaNumeric(b byte) bool {
	return isAlpha(b) || (b >= '0' && b <= '9')
}
