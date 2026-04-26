// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1383",
		Title:    "Avoid `$TIMEFORMAT` — Zsh uses `$TIMEFMT`",
		Severity: SeverityWarning,
		Description: "Bash's `$TIMEFORMAT` controls the output of the `time` builtin. Zsh uses a " +
			"shorter name, `$TIMEFMT`, for the same purpose. Setting `TIMEFORMAT` in a Zsh script " +
			"has no effect; the Zsh `time` builtin reads `$TIMEFMT`.",
		Check: checkZC1383,
		Fix:   fixZC1383,
	})
}

// fixZC1383 renames every `TIMEFORMAT` token inside an echo / print /
// printf / export argument to `TIMEFMT`. Each occurrence becomes its own
// edit at the absolute source offset of that arg's token + the substring
// index; surrounding quoting and adjoining text stay byte-exact.
func fixZC1383(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" && ident.Value != "export" {
		return nil
	}
	var edits []FixEdit
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if !strings.Contains(val, "TIMEFORMAT") {
			continue
		}
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+len(val) > len(source) {
			continue
		}
		if string(source[off:off+len(val)]) != val {
			continue
		}
		idx := 0
		for {
			pos := strings.Index(val[idx:], "TIMEFORMAT")
			if pos < 0 {
				break
			}
			abs := off + idx + pos
			line, col := offsetLineColZC1383(source, abs)
			if line < 0 {
				break
			}
			edits = append(edits, FixEdit{
				Line:    line,
				Column:  col,
				Length:  len("TIMEFORMAT"),
				Replace: "TIMEFMT",
			})
			idx += pos + len("TIMEFORMAT")
		}
	}
	return edits
}

func offsetLineColZC1383(source []byte, offset int) (int, int) {
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

func checkZC1383(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" && ident.Value != "export" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "TIMEFORMAT") {
			return []Violation{{
				KataID: "ZC1383",
				Message: "`$TIMEFORMAT` is Bash-only. Zsh reads `$TIMEFMT` (shorter name) for the " +
					"`time` builtin's output format.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
