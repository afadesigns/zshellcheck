// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1381",
		Title:    "Avoid `$COMP_WORDS`/`$COMP_CWORD` — Zsh uses `words`/`$CURRENT`",
		Severity: SeverityError,
		Description: "Bash programmable completion reads the partial command via `$COMP_WORDS` " +
			"(array of tokens) and `$COMP_CWORD` (index of cursor). Zsh's completion system " +
			"exposes the same via `words` (array) and `$CURRENT` (1-based cursor index). Using " +
			"the Bash names in Zsh completion functions produces empty expansions.",
		Check: checkZC1381,
		Fix:   fixZC1381,
	})
}

// fixZC1381 rewrites Bash completion variable names inside echo /
// print / printf args to their Zsh equivalents:
//
//	COMP_WORDS  → words
//	COMP_CWORD  → CURRENT
//	COMP_LINE   → BUFFER
//	COMP_POINT  → CURSOR
//
// Per-arg byte-anchored scan; one edit per match. Idempotent — a
// re-run sees the Zsh names, which the detector's substring guard
// won't match.
func fixZC1381(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" {
		return nil
	}
	mapping := []struct{ old, new string }{
		{"COMP_WORDS", "words"},
		{"COMP_CWORD", "CURRENT"},
		{"COMP_LINE", "BUFFER"},
		{"COMP_POINT", "CURSOR"},
	}
	var edits []FixEdit
	for _, arg := range cmd.Arguments {
		val := arg.String()
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+len(val) > len(source) {
			continue
		}
		if string(source[off:off+len(val)]) != val {
			continue
		}
		for _, m := range mapping {
			idx := 0
			for {
				pos := strings.Index(val[idx:], m.old)
				if pos < 0 {
					break
				}
				abs := off + idx + pos
				line, col := offsetLineColZC1381(source, abs)
				if line < 0 {
					break
				}
				edits = append(edits, FixEdit{
					Line:    line,
					Column:  col,
					Length:  len(m.old),
					Replace: m.new,
				})
				idx += pos + len(m.old)
			}
		}
	}
	return edits
}

func offsetLineColZC1381(source []byte, offset int) (int, int) {
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

func checkZC1381(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "COMP_WORDS") || strings.Contains(v, "COMP_CWORD") ||
			strings.Contains(v, "COMP_LINE") || strings.Contains(v, "COMP_POINT") {
			return []Violation{{
				KataID: "ZC1381",
				Message: "Bash `$COMP_*` completion variables do not exist in Zsh. Use " +
					"`$words` (array of tokens), `$CURRENT` (cursor index), `$BUFFER`, or the " +
					"`_arguments`/`_values` helpers from `compsys`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
