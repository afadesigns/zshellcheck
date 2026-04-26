// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1377",
		Title:    "Avoid `$BASH_ALIASES` — use Zsh `$aliases` associative array",
		Severity: SeverityWarning,
		Description: "Bash's `$BASH_ALIASES` is an associative array of alias→value mappings. Zsh " +
			"exposes the same information via `$aliases` (also an assoc array). `$BASH_ALIASES` " +
			"is unset in Zsh; reading it yields nothing.",
		Check: checkZC1377,
		Fix:   fixZC1377,
	})
}

// fixZC1377 renames every `BASH_ALIASES` token inside an echo / print /
// printf argument to `aliases`. Each occurrence becomes its own edit at
// the absolute source offset of that arg's token + the substring index;
// surrounding quoting and adjoining text stay byte-exact.
func fixZC1377(node ast.Node, _ Violation, source []byte) []FixEdit {
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
	var edits []FixEdit
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if !strings.Contains(val, "BASH_ALIASES") {
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
			pos := strings.Index(val[idx:], "BASH_ALIASES")
			if pos < 0 {
				break
			}
			abs := off + idx + pos
			line, col := offsetLineColZC1377(source, abs)
			if line < 0 {
				break
			}
			edits = append(edits, FixEdit{
				Line:    line,
				Column:  col,
				Length:  len("BASH_ALIASES"),
				Replace: "aliases",
			})
			idx += pos + len("BASH_ALIASES")
		}
	}
	return edits
}

func offsetLineColZC1377(source []byte, offset int) (int, int) {
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

func checkZC1377(node ast.Node) []Violation {
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
		if strings.Contains(v, "BASH_ALIASES") {
			return []Violation{{
				KataID: "ZC1377",
				Message: "`$BASH_ALIASES` is Bash-only. In Zsh use `$aliases` (assoc array) — " +
					"same structure, e.g. `print -l ${(kv)aliases}`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
