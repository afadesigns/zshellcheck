// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1238",
		Title:    "Avoid `docker exec -it` in scripts — drop `-it` for non-interactive",
		Severity: SeverityWarning,
		Description: "`docker exec -it` allocates a TTY and attaches stdin, which hangs " +
			"in non-interactive scripts. Use `docker exec` without `-it` for scripted commands.",
		Check: checkZC1238,
		Fix:   fixZC1238,
	})
}

// fixZC1238 strips the `-it` (or `-ti`) flag from a `docker exec`
// invocation. The span covers the leading whitespace plus the flag
// token so the surrounding source stays byte-identical.
func fixZC1238(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}
	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "exec" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		if v := arg.String(); v == "-it" || v == "-ti" {
			return zc1238StripFlag(source, arg, v)
		}
	}
	return nil
}

// zc1238StripFlag deletes the flag arg plus the run of horizontal
// whitespace immediately preceding it; the leading space the user
// typed disappears with the flag, leaving `docker exec CMD`.
func zc1238StripFlag(source []byte, arg ast.Expression, lit string) []FixEdit {
	tok := arg.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+len(lit) > len(source) {
		return nil
	}
	if string(source[off:off+len(lit)]) != lit {
		return nil
	}
	start := off
	for start > 0 && (source[start-1] == ' ' || source[start-1] == '\t') {
		start--
	}
	end := off + len(lit)
	startLine, startCol := offsetLineColZC1238(source, start)
	if startLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    startLine,
		Column:  startCol,
		Length:  end - start,
		Replace: "",
	}}
}

func offsetLineColZC1238(source []byte, offset int) (int, int) {
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

func checkZC1238(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "exec" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "-it" || val == "-ti" {
			return []Violation{{
				KataID: "ZC1238",
				Message: "Avoid `docker exec -it` in scripts — TTY allocation hangs without a terminal. " +
					"Use `docker exec` without `-it` for non-interactive commands.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
