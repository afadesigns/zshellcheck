// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:       "ZC1190",
		Title:    "Combine chained `grep -v` into single invocation",
		Severity: SeverityStyle,
		Description: "`grep -v p1 | grep -v p2` spawns two processes. " +
			"Use `grep -v -e p1 -e p2` to combine exclusions in one invocation.",
		Check: checkZC1190,
		Fix:   fixZC1190,
	})
}

// fixZC1190 collapses `grep -v p1 | grep -v p2` into a single
// `grep -v -e p1 -e p2`. Only fires when each grep has exactly one
// non-flag pattern argument and at most the lone `-v` flag — keeps the
// rewrite safe in the presence of trailing FILE / additional flags.
func fixZC1190(node ast.Node, _ Violation, source []byte) []FixEdit {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}
	left, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(left, "grep") {
		return nil
	}
	right, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok || !isCommandName(right, "grep") {
		return nil
	}
	leftPat, leftOk := zc1190SinglePattern(left)
	rightPat, rightOk := zc1190SinglePattern(right)
	if !leftOk || !rightOk {
		return nil
	}

	leftTok := left.TokenLiteralNode()
	leftStart := LineColToByteOffset(source, leftTok.Line, leftTok.Column)
	if leftStart < 0 {
		return nil
	}
	if len(right.Arguments) == 0 {
		return nil
	}
	lastArg := right.Arguments[len(right.Arguments)-1]
	laTok := lastArg.TokenLiteralNode()
	laOff := LineColToByteOffset(source, laTok.Line, laTok.Column)
	if laOff < 0 {
		return nil
	}
	laLit := lastArg.String()
	end := laOff + len(laLit)
	if end > len(source) {
		return nil
	}
	startLine, startCol := offsetLineColZC1190(source, leftStart)
	if startLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    startLine,
		Column:  startCol,
		Length:  end - leftStart,
		Replace: "grep -v -e " + leftPat + " -e " + rightPat,
	}}
}

// zc1190SinglePattern returns the lone non-flag argument of a
// `grep -v PAT` invocation. Returns ok=false when args contain
// extras, multiple flags, or no pattern at all.
func zc1190SinglePattern(cmd *ast.SimpleCommand) (string, bool) {
	pattern := ""
	hasV := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch {
		case v == "-v":
			hasV = true
		case len(v) > 0 && v[0] == '-':
			return "", false
		default:
			if pattern != "" {
				return "", false
			}
			pattern = v
		}
	}
	if !hasV || pattern == "" {
		return "", false
	}
	return pattern, true
}

func offsetLineColZC1190(source []byte, offset int) (int, int) {
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

func checkZC1190(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	leftCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(leftCmd, "grep") {
		return nil
	}

	rightCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok || !isCommandName(rightCmd, "grep") {
		return nil
	}

	leftHasV := false
	rightHasV := false

	for _, arg := range leftCmd.Arguments {
		if arg.String() == "-v" {
			leftHasV = true
		}
	}
	for _, arg := range rightCmd.Arguments {
		if arg.String() == "-v" {
			rightHasV = true
		}
	}

	if leftHasV && rightHasV {
		return []Violation{{
			KataID: "ZC1190",
			Message: "Combine `grep -v p1 | grep -v p2` into `grep -v -e p1 -e p2`. " +
				"A single invocation avoids an unnecessary pipeline.",
			Line:   pipe.TokenLiteralNode().Line,
			Column: pipe.TokenLiteralNode().Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
