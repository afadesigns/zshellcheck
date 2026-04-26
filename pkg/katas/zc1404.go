// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1404",
		Title:    "Avoid `$BASH_CMDS` — Bash-specific hash-table mirror, use Zsh `$commands`",
		Severity: SeverityWarning,
		Description: "Bash's `$BASH_CMDS` associative array mirrors the hash-table of command " +
			"names→paths. Zsh exposes the same via `$commands` (assoc array from " +
			"`zsh/parameter`). `$BASH_CMDS` is unset in Zsh.",
		Check: checkZC1404,
		Fix:   fixZC1404,
	})
}

// fixZC1404 rewrites `BASH_CMDS` → `commands` inside echo / print /
// printf args. Per-arg substring scan; one edit per match.
// Idempotent — a re-run sees `commands`, which the detector's
// substring guard won't match.
func fixZC1404(node ast.Node, _ Violation, source []byte) []FixEdit {
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
	const oldName = "BASH_CMDS"
	const newName = "commands"
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
		idx := 0
		for {
			pos := strings.Index(val[idx:], oldName)
			if pos < 0 {
				break
			}
			abs := off + idx + pos
			line, col := offsetLineColZC1404(source, abs)
			if line < 0 {
				break
			}
			edits = append(edits, FixEdit{
				Line:    line,
				Column:  col,
				Length:  len(oldName),
				Replace: newName,
			})
			idx += pos + len(oldName)
		}
	}
	return edits
}

func offsetLineColZC1404(source []byte, offset int) (int, int) {
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

func checkZC1404(node ast.Node) []Violation {
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
		if strings.Contains(v, "BASH_CMDS") {
			return []Violation{{
				KataID: "ZC1404",
				Message: "`$BASH_CMDS` is Bash-only. In Zsh use `$commands` (assoc array, " +
					"names→paths) via `zsh/parameter`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
