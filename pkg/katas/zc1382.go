package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1382",
		Title:    "Avoid `$READLINE_LINE`/`$READLINE_POINT` — Zsh ZLE uses `$BUFFER`/`$CURSOR`",
		Severity: SeverityError,
		Description: "Bash readline exposes the current input line as `$READLINE_LINE` and cursor " +
			"offset as `$READLINE_POINT` inside `bind -x` handlers. Zsh's Line Editor (ZLE) uses " +
			"`$BUFFER` (line text) and `$CURSOR` (1-based column) inside widget functions. The " +
			"Bash names are unset in Zsh.",
		Check: checkZC1382,
		Fix:   fixZC1382,
	})
}

// fixZC1382 rewrites Bash readline variable names inside echo /
// print / printf args to their Zsh ZLE equivalents:
//
//	READLINE_LINE   → BUFFER
//	READLINE_POINT  → CURSOR
//	READLINE_MARK   → MARK
//
// Per-arg byte-anchored scan; one edit per match. Idempotent — a
// re-run sees the Zsh names, which the detector's substring guard
// won't match.
func fixZC1382(node ast.Node, _ Violation, source []byte) []FixEdit {
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
		{"READLINE_LINE", "BUFFER"},
		{"READLINE_POINT", "CURSOR"},
		{"READLINE_MARK", "MARK"},
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
				line, col := offsetLineColZC1382(source, abs)
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

func offsetLineColZC1382(source []byte, offset int) (int, int) {
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

func checkZC1382(node ast.Node) []Violation {
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
		if strings.Contains(v, "READLINE_LINE") || strings.Contains(v, "READLINE_POINT") ||
			strings.Contains(v, "READLINE_MARK") {
			return []Violation{{
				KataID: "ZC1382",
				Message: "Bash `$READLINE_*` vars do not exist in Zsh. Inside ZLE widgets use " +
					"`$BUFFER`, `$CURSOR`, `$MARK`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
