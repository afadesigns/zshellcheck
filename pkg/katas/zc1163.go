package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:       "ZC1163",
		Title:    "Use `grep -m 1` instead of `grep | head -1`",
		Severity: SeverityStyle,
		Description: "`grep pattern | head -1` spawns two processes when `grep -m 1` does the same. " +
			"The `-m` flag stops after the first match, avoiding the pipeline.",
		Check: checkZC1163,
		Fix:   fixZC1163,
	})
}

// fixZC1163 collapses `grep PAT | head -1` into `grep -m 1 PAT`. Span
// runs from just after the `grep` command name through the end of the
// `head -1` invocation; the replacement preserves every original grep
// argument verbatim and drops the pipe + head suffix in one edit. Only
// fires for the `-1` / `-n1` shapes the detector already guards.
func fixZC1163(node ast.Node, _ Violation, source []byte) []FixEdit {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}
	grepCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(grepCmd, "grep") {
		return nil
	}
	headCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok || !isCommandName(headCmd, "head") {
		return nil
	}
	hasFirst := false
	for _, arg := range headCmd.Arguments {
		v := arg.String()
		if v == "-1" || v == "-n1" {
			hasFirst = true
		}
	}
	if !hasFirst {
		return nil
	}

	grepTok := grepCmd.TokenLiteralNode()
	grepOff := LineColToByteOffset(source, grepTok.Line, grepTok.Column)
	if grepOff < 0 {
		return nil
	}
	grepLen := IdentLenAt(source, grepOff)
	if grepLen == 0 {
		return nil
	}
	spanStart := grepOff + grepLen

	pipeOff := LineColToByteOffset(source, pipe.Token.Line, pipe.Token.Column)
	if pipeOff < 0 || pipeOff >= len(source) || source[pipeOff] != '|' {
		return nil
	}
	argsEnd := pipeOff
	for argsEnd > spanStart && (source[argsEnd-1] == ' ' || source[argsEnd-1] == '\t') {
		argsEnd--
	}
	middle := string(source[spanStart:argsEnd])

	// End of head -1: the last argument's last byte.
	if len(headCmd.Arguments) == 0 {
		return nil
	}
	lastArg := headCmd.Arguments[len(headCmd.Arguments)-1]
	lastTok := lastArg.TokenLiteralNode()
	lastOff := LineColToByteOffset(source, lastTok.Line, lastTok.Column)
	if lastOff < 0 {
		return nil
	}
	lastLit := lastArg.String()
	spanEnd := lastOff + len(lastLit)
	if spanEnd <= spanStart {
		return nil
	}

	startLine, startCol := offsetLineColZC1163(source, spanStart)
	if startLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    startLine,
		Column:  startCol,
		Length:  spanEnd - spanStart,
		Replace: " -m 1" + middle,
	}}
}

func offsetLineColZC1163(source []byte, offset int) (int, int) {
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

func checkZC1163(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	grepCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(grepCmd, "grep") {
		return nil
	}

	headCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok || !isCommandName(headCmd, "head") {
		return nil
	}

	// Check head has -1 or -n 1
	for _, arg := range headCmd.Arguments {
		val := arg.String()
		if val == "-1" || val == "-n1" {
			return []Violation{{
				KataID: "ZC1163",
				Message: "Use `grep -m 1` instead of `grep | head -1`. " +
					"The `-m` flag stops after the first match without a pipeline.",
				Line:   pipe.TokenLiteralNode().Line,
				Column: pipe.TokenLiteralNode().Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
