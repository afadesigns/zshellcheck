// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
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
	_, catCmd, rightCmd, _, ok := zc1146Pipe(node)
	if !ok {
		return nil
	}
	catStart, ok := zc1146Offset(source, catCmd.TokenLiteralNode())
	if !ok {
		return nil
	}
	fileLit, _, ok := zc1146ArgSlice(source, catCmd.Arguments[0])
	if !ok {
		return nil
	}
	rightStart, ok := zc1146Offset(source, rightCmd.TokenLiteralNode())
	if !ok {
		return nil
	}
	rightEnd, ok := zc1146RightEnd(source, rightCmd, rightStart)
	if !ok {
		return nil
	}
	startLine, startCol := offsetLineColZC1146(source, catStart)
	if startLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    startLine,
		Column:  startCol,
		Length:  rightEnd - catStart,
		Replace: string(source[rightStart:rightEnd]) + " " + fileLit,
	}}
}

func zc1146Offset(source []byte, tok token.Token) (int, bool) {
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	return off, off >= 0
}

// zc1146ArgSlice returns the literal text of arg as it appears in
// source plus the offset, or ok=false when the AST coordinates do not
// line up with the source bytes.
func zc1146ArgSlice(source []byte, arg ast.Expression) (lit string, off int, ok bool) {
	tok := arg.TokenLiteralNode()
	off, ok = zc1146Offset(source, tok)
	if !ok {
		return "", 0, false
	}
	lit = arg.String()
	if off+len(lit) > len(source) || string(source[off:off+len(lit)]) != lit {
		return "", 0, false
	}
	return lit, off, true
}

func zc1146RightEnd(source []byte, rightCmd *ast.SimpleCommand, rightStart int) (int, bool) {
	rightIdent, ok := rightCmd.Name.(*ast.Identifier)
	if !ok {
		return 0, false
	}
	end := rightStart + len(rightIdent.Value)
	if n := len(rightCmd.Arguments); n > 0 {
		lastArg := rightCmd.Arguments[n-1]
		laOff, ok := zc1146Offset(source, lastArg.TokenLiteralNode())
		if !ok {
			return 0, false
		}
		end = laOff + len(lastArg.String())
	}
	if end > len(source) || end < rightStart {
		return 0, false
	}
	return end, true
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

var zc1146FileTakers = map[string]struct{}{
	"awk":  {},
	"sed":  {},
	"sort": {},
	"head": {},
	"tail": {},
}

func checkZC1146(node ast.Node) []Violation {
	pipe, _, _, name, ok := zc1146Pipe(node)
	if !ok {
		return nil
	}
	if _, hit := zc1146FileTakers[name]; !hit {
		return nil
	}
	return []Violation{{
		KataID: "ZC1146",
		Message: "Pass the file directly to `" + name + "` instead of `cat file | " + name + "`. " +
			"Most text-processing tools accept file arguments.",
		Line:   pipe.TokenLiteralNode().Line,
		Column: pipe.TokenLiteralNode().Column,
		Level:  SeverityStyle,
	}}
}

// zc1146Pipe destructures `cat FILE | NAME [args]` into its parts and
// reports whether the cat side is well-formed (single non-flag arg).
func zc1146Pipe(node ast.Node) (pipe *ast.InfixExpression, catCmd, rightCmd *ast.SimpleCommand, name string, ok bool) {
	pipe, isPipe := node.(*ast.InfixExpression)
	if !isPipe || pipe.Operator != "|" {
		return nil, nil, nil, "", false
	}
	catCmd, isCat := pipe.Left.(*ast.SimpleCommand)
	if !isCat || !isCommandName(catCmd, "cat") || len(catCmd.Arguments) != 1 {
		return nil, nil, nil, "", false
	}
	if first := catCmd.Arguments[0].String(); first != "" && first[0] == '-' {
		return nil, nil, nil, "", false
	}
	rightCmd, isRight := pipe.Right.(*ast.SimpleCommand)
	if !isRight {
		return nil, nil, nil, "", false
	}
	rightIdent, isIdent := rightCmd.Name.(*ast.Identifier)
	if !isIdent {
		return nil, nil, nil, "", false
	}
	return pipe, catCmd, rightCmd, rightIdent.Value, true
}
