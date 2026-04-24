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
func fixZC1091(node ast.Node, v Violation, source []byte) []FixEdit {
	dbe, ok := node.(*ast.DoubleBracketExpression)
	if !ok {
		return nil
	}
	openLine := dbe.Token.Line
	openCol := dbe.Token.Column
	openOff := LineColToByteOffset(source, openLine, openCol)
	if openOff < 0 {
		return nil
	}
	// The lexer stamps `[[` at the second bracket (two-char fusion).
	// Step back one column when needed so the edit covers the whole
	// opener.
	if openOff > 0 && source[openOff] == '[' && source[openOff-1] == '[' {
		openOff--
		openCol--
	}
	if openOff+2 > len(source) || source[openOff] != '[' || source[openOff+1] != '[' {
		return nil
	}
	closeOff := findDoubleBracketClose(source, openOff+2)
	if closeOff < 0 {
		return nil
	}
	// Find the single `-eq`/etc. operator token and replace.
	var opTok ast.Expression
	var opStr string
	opCount := 0
	for _, el := range dbe.Elements {
		if infix, ok := el.(*ast.InfixExpression); ok {
			if repl, found := arithCmpReplacements[infix.Operator]; found {
				opTok = el
				opStr = infix.Operator
				_ = repl
				opCount++
			}
		}
	}
	if opCount != 1 || opTok == nil {
		return nil
	}
	// opTok points at the infix expression; its TokenLiteralNode is
	// the operator token (e.g. `-eq`). Row/col is already 1-based.
	infix := opTok.(*ast.InfixExpression)
	opLine := infix.Token.Line
	opCol := infix.Token.Column
	closeLine, closeCol := offsetLineColZC1091(source, closeOff)
	if closeLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: openLine, Column: openCol, Length: 2, Replace: "(("},
		{Line: opLine, Column: opCol, Length: len(opStr), Replace: arithCmpReplacements[opStr]},
		{Line: closeLine, Column: closeCol, Length: 2, Replace: "))"},
	}
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
