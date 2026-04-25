package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:       "ZC1146",
		Title:    "Avoid `cat file | awk` — pass file to awk directly",
		Severity: SeverityStyle,
		Description: "`cat file | awk` spawns an unnecessary cat process. " +
			"Pass the file directly as `awk '...' file`.",
		Check: checkZC1146,
		Fix:   fixZC1146,
	})
}

// fixZC1146 collapses `cat FILE | tool [args]` into `tool [args] FILE`.
// One span replacement runs from the start of `cat` through the end of
// the right-hand command; the replacement is the right-hand source
// verbatim with ` FILE` appended. Only fires when the cat command has
// exactly one filename argument (the detector already guards that).
func fixZC1146(node ast.Node, _ Violation, source []byte) []FixEdit {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}
	catCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(catCmd, "cat") {
		return nil
	}
	if len(catCmd.Arguments) != 1 {
		return nil
	}
	rightCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	rightIdent, ok := rightCmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	catTok := catCmd.TokenLiteralNode()
	catStart := LineColToByteOffset(source, catTok.Line, catTok.Column)
	if catStart < 0 {
		return nil
	}

	// File arg literal text: take it byte-exact from the source.
	fileArg := catCmd.Arguments[0]
	fileTok := fileArg.TokenLiteralNode()
	fileOff := LineColToByteOffset(source, fileTok.Line, fileTok.Column)
	if fileOff < 0 {
		return nil
	}
	fileLit := fileArg.String()
	if fileOff+len(fileLit) > len(source) ||
		string(source[fileOff:fileOff+len(fileLit)]) != fileLit {
		return nil
	}

	// Right command source span: [rightStart, rightEnd).
	rightTok := rightCmd.TokenLiteralNode()
	rightStart := LineColToByteOffset(source, rightTok.Line, rightTok.Column)
	if rightStart < 0 {
		return nil
	}
	rightEnd := rightStart + len(rightIdent.Value)
	if len(rightCmd.Arguments) > 0 {
		lastArg := rightCmd.Arguments[len(rightCmd.Arguments)-1]
		laTok := lastArg.TokenLiteralNode()
		laOff := LineColToByteOffset(source, laTok.Line, laTok.Column)
		if laOff < 0 {
			return nil
		}
		laLit := lastArg.String()
		rightEnd = laOff + len(laLit)
	}
	if rightEnd > len(source) || rightEnd < rightStart {
		return nil
	}
	rightSrc := string(source[rightStart:rightEnd])

	startLine, startCol := offsetLineColZC1146(source, catStart)
	if startLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    startLine,
		Column:  startCol,
		Length:  rightEnd - catStart,
		Replace: rightSrc + " " + fileLit,
	}}
}

func offsetLineColZC1146(source []byte, offset int) (int, int) {
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

func checkZC1146(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	catCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(catCmd, "cat") {
		return nil
	}

	// cat must have exactly one file argument and no flags
	if len(catCmd.Arguments) != 1 {
		return nil
	}
	for _, arg := range catCmd.Arguments {
		if len(arg.String()) > 0 && arg.String()[0] == '-' {
			return nil
		}
	}

	// Right side must be awk/sed/sort or similar
	rightCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	rightIdent, ok := rightCmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := rightIdent.Value
	if name == "awk" || name == "sed" || name == "sort" || name == "head" || name == "tail" {
		return []Violation{{
			KataID: "ZC1146",
			Message: "Pass the file directly to `" + name + "` instead of `cat file | " + name + "`. " +
				"Most text-processing tools accept file arguments.",
			Line:   pipe.TokenLiteralNode().Line,
			Column: pipe.TokenLiteralNode().Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
