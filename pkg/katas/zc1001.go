// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IndexExpressionNode, Kata{
		ID:    "ZC1001",
		Title: "Use ${} for array element access",
		Description: "In Zsh, accessing array elements with `$my_array[1]` doesn't work as expected. " +
			"It tries to access an element from an array named `my_array[1]`. " +
			"The correct way to access an array element is to use `${my_array[1]}`.",
		Severity: SeverityStyle,
		Check:    checkZC1001,
		Fix:      fixZC1001,
	})
	RegisterKata(ast.InvalidArrayAccessNode, Kata{
		ID:    "ZC1001",
		Title: "Use ${} for array element access",
		Description: "In Zsh, accessing array elements with `$my_array[1]` doesn't work as expected. " +
			"It tries to access an element from an array named `my_array[1]`. " +
			"The correct way to access an array element is to use `${my_array[1]}`.",
		Severity: SeverityStyle,
		Check:    checkZC1001,
		Fix:      fixZC1001,
	})
}

// fixZC1001 rewrites `$arr[i]` to `${arr[i]}`. Two edits: insert
// `{` between the `$` and the identifier, then insert `}` after
// the closing `]`. Source positions are derived from the violation
// column (which points at the leading `$`) and a quote/brace-aware
// scan for the matching `]`.
func fixZC1001(node ast.Node, v Violation, source []byte) []FixEdit {
	dollarOff := LineColToByteOffset(source, v.Line, v.Column)
	if dollarOff < 0 || dollarOff >= len(source) || source[dollarOff] != '$' {
		return nil
	}
	// Find the `[` that opens the subscript starting from after `$name`.
	// Walk identifier chars then expect `[`.
	i := dollarOff + 1
	for i < len(source) && (isIdentByte(source[i])) {
		i++
	}
	if i >= len(source) || source[i] != '[' {
		return nil
	}
	closeOff := findSubscriptClose(source, i)
	if closeOff < 0 {
		return nil
	}
	return []FixEdit{
		// Insert `{` immediately after the `$`.
		{Line: v.Line, Column: v.Column + 1, Length: 0, Replace: "{"},
		// Insert `}` immediately after the closing `]`.
		offsetToEdit(source, closeOff+1, 0, "}"),
	}
}

// findSubscriptClose returns the byte offset of the `]` that closes
// the subscript opened at open. Tracks single / double quotes and
// nested `[` so subscripts containing patterns like `arr[$other[i]]`
// or `arr[(R)pat]` close at the right bracket. Returns -1 when no
// match is found before EOF.
func findSubscriptClose(source []byte, open int) int {
	inSingle := false
	inDouble := false
	depth := 1
	for i := open + 1; i < len(source); i++ {
		c := source[i]
		switch {
		case inSingle:
			if c == '\'' {
				inSingle = false
			}
		case inDouble:
			if c == '\\' && i+1 < len(source) {
				i++
				continue
			}
			if c == '"' {
				inDouble = false
			}
		default:
			switch c {
			case '\'':
				inSingle = true
			case '"':
				inDouble = true
			case '[':
				depth++
			case ']':
				depth--
				if depth == 0 {
					return i
				}
			case '\n':
				return -1
			}
		}
	}
	return -1
}

func checkZC1001(node ast.Node) []Violation {
	violations := []Violation{}

	if indexExp, ok := node.(*ast.IndexExpression); ok {
		if ident, ok := indexExp.Left.(*ast.Identifier); ok {
			if len(ident.Value) > 0 && ident.Value[0] == '$' {
				violations = append(violations, Violation{
					KataID: "ZC1001",
					Message: "Use ${} for array element access. " +
						"Accessing array elements with `" + ident.Value + "[...]` is not the correct syntax in Zsh.",
					Line:   ident.Token.Line,
					Column: ident.Token.Column,
					Level:  SeverityStyle,
				})
			}
		}
	} else if arrayAccess, ok := node.(*ast.InvalidArrayAccess); ok {
		violations = append(violations, Violation{
			KataID: "ZC1001",
			Message: "Use ${} for array element access. " +
				"Accessing array elements with `$my_array[1]` is not the correct syntax in Zsh.",
			Line:   arrayAccess.Token.Line,
			Column: arrayAccess.Token.Column,
			Level:  SeverityStyle,
		})
	}

	return violations
}
