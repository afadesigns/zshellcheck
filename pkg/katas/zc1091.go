// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.DoubleBracketExpressionNode, Kata{
		ID:    "ZC1091",
		Title: "Use `((...))` for arithmetic comparisons in `[[...]]`",
		Description: "The `[[ ... ]]` construct is primarily for string comparisons and file tests. " +
			"For arithmetic comparisons (`-eq`, `-lt`, etc.), use the dedicated arithmetic context `(( ... ))`. " +
			"It is cleaner and strictly numeric.",
		Severity: SeverityStyle,
		Check:    checkZC1091,
		Fix:      fixZC1091,
	})
}

// fixZC1091 rewrites a bracket conditional that uses dashed
// comparison operators into arithmetic form. Example:
// `[[ x -lt 10 ]]` → `(( x < 10 ))`. Only fires when exactly one
// recognised operator appears inside the brackets to keep the
// rewrite unambiguous.
func fixZC1091(node ast.Node, _ Violation, source []byte) []FixEdit {
	dbe, ok := node.(*ast.DoubleBracketExpression)
	if !ok {
		return nil
	}
	openOff, openLine, openCol, ok := zc1091OpenBracket(source, dbe)
	if !ok {
		return nil
	}
	closeOff := findDoubleBracketClose(source, openOff+2)
	if closeOff < 0 {
		return nil
	}
	infix, ok := zc1091SingleArithOp(dbe)
	if !ok {
		return nil
	}
	closeLine, closeCol := offsetLineColZC1091(source, closeOff)
	if closeLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: openLine, Column: openCol, Length: 2, Replace: "(("},
		{Line: infix.Token.Line, Column: infix.Token.Column, Length: len(infix.Operator), Replace: arithCmpReplacements[infix.Operator]},
		{Line: closeLine, Column: closeCol, Length: 2, Replace: "))"},
	}
}

func zc1091OpenBracket(source []byte, dbe *ast.DoubleBracketExpression) (off, line, col int, ok bool) {
	line = dbe.Token.Line
	col = dbe.Token.Column
	off = LineColToByteOffset(source, line, col)
	if off < 0 {
		return 0, 0, 0, false
	}
	if off > 0 && source[off] == '[' && source[off-1] == '[' {
		off--
		col--
	}
	if off+2 > len(source) || source[off] != '[' || source[off+1] != '[' {
		return 0, 0, 0, false
	}
	return off, line, col, true
}

func zc1091SingleArithOp(dbe *ast.DoubleBracketExpression) (*ast.InfixExpression, bool) {
	var found *ast.InfixExpression
	for _, el := range dbe.Elements {
		infix, ok := el.(*ast.InfixExpression)
		if !ok {
			continue
		}
		if _, hit := arithCmpReplacements[infix.Operator]; !hit {
			continue
		}
		if found != nil {
			return nil, false
		}
		found = infix
	}
	return found, found != nil
}

// findDoubleBracketClose scans source for the matching `]]` that
// closes the `[[` just before `start`. Honours `[…]` nesting so
// character classes like `[:alnum:]` don't trip the scan.
func findDoubleBracketClose(source []byte, start int) int {
	depth := 0
	for i := start; i < len(source)-1; i++ {
		switch source[i] {
		case '\\':
			i++
		case '[':
			depth++
		case ']':
			if depth > 0 {
				depth--
				continue
			}
			if source[i+1] == ']' {
				return i
			}
		}
	}
	return -1
}

func offsetLineColZC1091(source []byte, offset int) (int, int) {
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

func checkZC1091(node ast.Node) []Violation {
	dbe, ok := node.(*ast.DoubleBracketExpression)
	if !ok {
		return nil
	}

	var violations []Violation

	visitor := func(n ast.Node) bool {
		if infix, ok := n.(*ast.InfixExpression); ok {
			switch infix.Operator {
			case "-eq", "-ne", "-lt", "-le", "-gt", "-ge":
				violations = append(violations, Violation{
					KataID:  "ZC1091",
					Message: "Use `(( ... ))` for arithmetic comparisons. e.g. `(( a < b ))` instead of `[[ a -lt b ]]`.",
					Line:    infix.TokenLiteralNode().Line,
					Column:  infix.TokenLiteralNode().Column,
					Level:   SeverityStyle,
				})
			}
		}
		return true
	}

	for _, expr := range dbe.Elements {
		ast.Walk(expr, visitor)
	}

	return violations
}
