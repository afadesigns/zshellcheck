// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1055",
		Title: "Use `[[ -n/-z ]]` for empty string checks",
		Description: "Comparing with empty string is less idiomatic than using `[[ -z $var ]]` (is empty) " +
			"or `[[ -n $var ]]` (is not empty).",
		Severity: SeverityStyle,
		Check:    checkZC1055,
		Fix:      fixZC1055,
	})
}

func checkZC1055(node ast.Node) []Violation {
	expr, ok := node.(*ast.InfixExpression)
	if !ok {
		return nil
	}

	// Check for == "" or != ""
	if expr.Operator != "==" && expr.Operator != "!=" {
		return nil
	}

	// Check if either side is empty string literal
	isEmptyString := func(n ast.Node) bool {
		if str, ok := n.(*ast.StringLiteral); ok {
			// Check for "" or ''
			val := str.Value
			return val == `""` || val == `''`
		}
		return false
	}

	if isEmptyString(expr.Left) || isEmptyString(expr.Right) {
		opSuggestion := "-z"
		if expr.Operator == "!=" {
			opSuggestion = "-n"
		}

		return []Violation{{
			KataID:  "ZC1055",
			Message: "Use `[[ " + opSuggestion + " ... ]]` instead of comparing with empty string.",
			Line:    expr.TokenLiteralNode().Line,
			Column:  expr.TokenLiteralNode().Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

// fixZC1055 rewrites `$var == ""` / `$var != ""` into `-z $var` /
// `-n $var` respectively. The span covers the full infix expression
// so `[[ $var == "" ]]` ends up as `[[ -z $var ]]`. Handles both
// left-side and right-side empty-string positions.
func fixZC1055(node ast.Node, v Violation, source []byte) []FixEdit {
	expr, ok := node.(*ast.InfixExpression)
	if !ok {
		return nil
	}
	if expr.Operator != "==" && expr.Operator != "!=" {
		return nil
	}
	isEmpty := func(n ast.Node) (bool, int) {
		str, ok := n.(*ast.StringLiteral)
		if !ok {
			return false, 0
		}
		if str.Value == `""` || str.Value == `''` {
			return true, len(str.Value)
		}
		return false, 0
	}

	var varNode ast.Node
	var emptyNode ast.Node
	var emptyLen int
	if ok, n := isEmpty(expr.Left); ok {
		varNode = expr.Right
		emptyNode = expr.Left
		emptyLen = n
	} else if ok, n := isEmpty(expr.Right); ok {
		varNode = expr.Left
		emptyNode = expr.Right
		emptyLen = n
	} else {
		return nil
	}

	varExpr, vok := varNode.(ast.Expression)
	emptyExpr, eok := emptyNode.(ast.Expression)
	if !vok || !eok {
		return nil
	}
	var varTok, emptyTok token.Token = varExpr.TokenLiteralNode(), emptyExpr.TokenLiteralNode()
	if varTok.Line == 0 || emptyTok.Line == 0 {
		return nil
	}
	varOffset := LineColToByteOffset(source, varTok.Line, varTok.Column)
	emptyOffset := LineColToByteOffset(source, emptyTok.Line, emptyTok.Column)
	if varOffset < 0 || emptyOffset < 0 {
		return nil
	}

	// Determine span start / end based on which side is the empty
	// literal. The span covers varNode + operator + emptyNode in
	// source order.
	start := varOffset
	if emptyOffset < start {
		start = emptyOffset
	}
	end := emptyOffset + emptyLen
	if varEnd := varOffset + identOrVarLen(source, varOffset); varEnd > end {
		end = varEnd
	}

	op := "-z"
	if expr.Operator == "!=" {
		op = "-n"
	}

	varText := string(source[varOffset : varOffset+identOrVarLen(source, varOffset)])
	line, col := v.Line, v.Column
	if startLine, startCol := byteOffsetToLineColZC1055(source, start); startLine > 0 {
		line, col = startLine, startCol
	}
	return []FixEdit{{
		Line:    line,
		Column:  col,
		Length:  end - start,
		Replace: op + " " + varText,
	}}
}

// identOrVarLen returns the byte length of an identifier or variable
// token that starts at offset. Variables may begin with `$`, `${`,
// or a plain identifier run. We scan until whitespace / delimiter so
// composite words like `$var.ext` stay together.
func identOrVarLen(source []byte, offset int) int {
	if offset < 0 || offset >= len(source) {
		return 0
	}
	n := 0
	depth := 0
	for offset+n < len(source) {
		c := source[offset+n]
		if c == ' ' || c == '\t' || c == '\n' {
			break
		}
		if depth == 0 {
			switch c {
			case ';', '|', '&', ')', ']', '}':
				return n
			}
		}
		if c == '{' || c == '(' {
			depth++
		} else if c == '}' || c == ')' {
			if depth > 0 {
				depth--
			}
		}
		n++
	}
	return n
}

func byteOffsetToLineColZC1055(source []byte, offset int) (int, int) {
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
