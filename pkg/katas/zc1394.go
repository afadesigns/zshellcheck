// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"regexp"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// bashVarRE matches `$BASH` used as a standalone variable (not `$BASH_`).
var bashVarRE = regexp.MustCompile(`\$BASH(?:[^_A-Z]|$)`)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1394",
		Title:    "Avoid `$BASH` — Zsh uses `$ZSH_NAME` for the interpreter name",
		Severity: SeverityInfo,
		Description: "Bash's `$BASH` holds the path to the running Bash executable. Zsh's " +
			"equivalent is `$ZSH_NAME` (for the binary name) or `$0` (interactive shell). " +
			"Using `$BASH` in a Zsh script yields empty output.",
		Check: checkZC1394,
		Fix:   fixZC1394,
	})
}

// fixZC1394 renames every `$BASH` token (not part of a longer
// `$BASH_*` identifier) inside an echo / print / printf argument to
// `$ZSH_NAME`. Each occurrence becomes its own edit at the absolute
// source offset of that arg's token + the substring index; surrounding
// quoting, trailing punctuation, and adjoining text stay byte-exact.
func fixZC1394(node ast.Node, _ Violation, source []byte) []FixEdit {
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
		matches := bashVarRE.FindAllStringIndex(val, -1)
		if len(matches) == 0 {
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
		for _, m := range matches {
			// The regex spans `$BASH` plus one trailing byte (or end).
			// Rewrite only the `$BASH` prefix; leave the trailing byte
			// (the boundary char such as a space or quote) intact.
			abs := off + m[0]
			line, col := offsetLineColZC1394(source, abs)
			if line < 0 {
				continue
			}
			edits = append(edits, FixEdit{
				Line:    line,
				Column:  col,
				Length:  len("$BASH"),
				Replace: "$ZSH_NAME",
			})
		}
	}
	return edits
}

func offsetLineColZC1394(source []byte, offset int) (int, int) {
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

func checkZC1394(node ast.Node) []Violation {
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
		if bashVarRE.MatchString(v) {
			return []Violation{{
				KataID: "ZC1394",
				Message: "`$BASH` is Bash-only. Zsh exposes the interpreter name via `$ZSH_NAME` " +
					"and the executable path indirectly via `$0`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}
