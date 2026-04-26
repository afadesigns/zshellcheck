// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1126",
		Title: "Use `sort -u` instead of `sort | uniq`",
		Description: "`sort | uniq` spawns two processes when `sort -u` does the same in one. " +
			"Use `sort -u` to deduplicate sorted output efficiently.",
		Severity: SeverityStyle,
		Check:    checkZC1126,
		Fix:      fixZC1126,
	})
}

// fixZC1126 collapses `sort ... | uniq` into `sort -u ...`. Uses a
// single span-replacement from just after the `sort` command name
// through the end of `uniq`, rewriting the region to ` -u` +
// whatever sort args sit between the name and the pipe. Only fires
// when `uniq` has no flags (ZC1126's detector already guards that).
func fixZC1126(node ast.Node, v Violation, source []byte) []FixEdit {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}
	sortCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	uniqCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	sortTok := sortCmd.TokenLiteralNode()
	sortNameOff := LineColToByteOffset(source, sortTok.Line, sortTok.Column)
	if sortNameOff < 0 {
		return nil
	}
	sortNameLen := IdentLenAt(source, sortNameOff)
	if sortNameLen == 0 {
		return nil
	}
	spanStart := sortNameOff + sortNameLen

	// Find the pipe byte and walk back past trailing whitespace.
	pipeOff := LineColToByteOffset(source, pipe.Token.Line, pipe.Token.Column)
	if pipeOff < 0 || source[pipeOff] != '|' {
		return nil
	}
	argsEnd := pipeOff
	for argsEnd > spanStart && (source[argsEnd-1] == ' ' || source[argsEnd-1] == '\t') {
		argsEnd--
	}
	middle := string(source[spanStart:argsEnd])

	// End of uniq: the identifier itself; detector forbids flags.
	uniqTok := uniqCmd.TokenLiteralNode()
	uniqOff := LineColToByteOffset(source, uniqTok.Line, uniqTok.Column)
	uniqLen := IdentLenAt(source, uniqOff)
	if uniqOff < 0 || uniqLen == 0 {
		return nil
	}
	spanEnd := uniqOff + uniqLen

	replace := " -u" + middle
	startLine, startCol := offsetLineColZC1126(source, spanStart)
	if startLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    startLine,
		Column:  startCol,
		Length:  spanEnd - spanStart,
		Replace: replace,
	}}
}

func offsetLineColZC1126(source []byte, offset int) (int, int) {
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

func checkZC1126(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	sortCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if !isCommandName(sortCmd, "sort") {
		return nil
	}

	uniqCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if !isCommandName(uniqCmd, "uniq") {
		return nil
	}

	// If uniq has flags like -c (count), -d (duplicates), skip
	for _, arg := range uniqCmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1126",
		Message: "Use `sort -u` instead of `sort | uniq`. " +
			"Combining into one command avoids an unnecessary pipeline.",
		Line:   pipe.TokenLiteralNode().Line,
		Column: pipe.TokenLiteralNode().Column,
		Level:  SeverityStyle,
	}}
}
