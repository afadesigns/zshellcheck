// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1643",
		Title:    "Style: `$(cat file)` — use `$(<file)` to skip the fork / exec",
		Severity: SeverityStyle,
		Description: "`$(cat FILE)` forks, execs `/usr/bin/cat`, reads FILE, writes the bytes " +
			"to the pipe, waits for the child. `$(<FILE)` is a shell builtin — it reads FILE " +
			"directly into the command-substitution buffer with no fork and no exec. In a hot " +
			"path the speedup is dramatic, and even in cold paths it avoids one of the most " +
			"common useless-use-of-cat patterns in review feedback.",
		Check: checkZC1643,
		Fix:   fixZC1643,
	})
}

// fixZC1643 rewrites `$(cat FILE)` into `$(<FILE)`. The detector
// matches on the literal `$(cat ` substring inside a command argument;
// the Fix scans each argument's source span for that prefix and
// replaces `cat ` with `<` (4-byte → 1-byte collapse). Each occurrence
// becomes its own FixEdit, so multiple `$(cat …)` substitutions on the
// same line are all handled in one pass.
//
// Idempotent because the literal `$(cat ` no longer appears in the
// rewritten text — the second-pass detector sees `$(<…)` and stays
// silent.
func fixZC1643(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	const needle = "$(cat "
	var edits []FixEdit
	for _, arg := range cmd.Arguments {
		tok := arg.TokenLiteralNode()
		argOff := LineColToByteOffset(source, tok.Line, tok.Column)
		if argOff < 0 {
			continue
		}
		// Constrain the search window to the argument's literal text
		// span. arg.String() may include AST decorations (parens for
		// expressions) so use the byte slice directly: walk forward
		// from the start of the argument until a delimiter that ends
		// the unquoted argument context.
		end := argOff
		// Track quoting so the search does not stop inside a quoted
		// substring. The argument span runs to the next unquoted
		// whitespace, `;`, `&`, `|`, `)` (when paren depth = 0), or
		// newline.
		inSingle := false
		inDouble := false
		parenDepth := 0
		for end < len(source) {
			c := source[end]
			if c == '\\' && end+1 < len(source) {
				end += 2
				continue
			}
			switch {
			case inSingle:
				if c == '\'' {
					inSingle = false
				}
			case inDouble:
				if c == '"' {
					inDouble = false
				}
			default:
				switch c {
				case '\'':
					inSingle = true
				case '"':
					inDouble = true
				case '(':
					parenDepth++
				case ')':
					if parenDepth == 0 {
						goto done
					}
					parenDepth--
				case ' ', '\t', '\n', ';', '&', '|':
					if parenDepth == 0 {
						goto done
					}
				}
			}
			end++
		}
	done:
		// Find every `$(cat ` occurrence inside [argOff, end).
		i := argOff
		for i+len(needle) <= end {
			if string(source[i:i+len(needle)]) != needle {
				i++
				continue
			}
			// Replace the `cat ` portion (4 bytes after `$(`) with `<`.
			line, col := offsetLineColZC1643(source, i+2)
			if line < 0 {
				i += len(needle)
				continue
			}
			edits = append(edits, FixEdit{
				Line:    line,
				Column:  col,
				Length:  4,
				Replace: "<",
			})
			i += len(needle)
		}
	}
	return edits
}

func offsetLineColZC1643(source []byte, offset int) (int, int) {
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

func checkZC1643(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "$(cat ") {
			return []Violation{{
				KataID: "ZC1643",
				Message: "`$(cat FILE)` forks cat just to read a file — use `$(<FILE)` " +
					"(shell builtin, no fork).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}
	return nil
}
