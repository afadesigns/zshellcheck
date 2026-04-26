// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1403",
		Title:    "Setting `$HISTFILESIZE` alone is incomplete in Zsh — pair with `$SAVEHIST`",
		Severity: SeverityWarning,
		Description: "Bash uses `$HISTSIZE` (in-memory) and `$HISTFILESIZE` (on disk). Zsh uses " +
			"`$HISTSIZE` (in-memory) and `$SAVEHIST` (on disk). Setting only `$HISTFILESIZE` in " +
			"Zsh has no effect on disk — `$SAVEHIST` must be set. Mixing both names leaves " +
			"disk-history behavior undefined.",
		Check: checkZC1403,
		Fix:   fixZC1403,
	})
}

// fixZC1403 rewrites `HISTFILESIZE` → `SAVEHIST` inside echo /
// print / printf / export args. Per-arg substring scan; one edit
// per match. Idempotent — a re-run sees `SAVEHIST`, which the
// detector's substring guard won't match.
func fixZC1403(node ast.Node, _ Violation, source []byte) []FixEdit {
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
	const oldName = "HISTFILESIZE"
	const newName = "SAVEHIST"
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
			line, col := offsetLineColZC1403(source, abs)
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

func offsetLineColZC1403(source []byte, offset int) (int, int) {
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

func checkZC1403(node ast.Node) []Violation {
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
		if strings.Contains(v, "HISTFILESIZE") {
			return []Violation{{
				KataID: "ZC1403",
				Message: "`$HISTFILESIZE` is Bash-only. Zsh uses `$SAVEHIST` for on-disk history " +
					"size. Setting `HISTFILESIZE` in Zsh has no effect.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
