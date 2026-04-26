// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1856",
		Title:    "Warn on `unset arr[N]` — Zsh does not delete the array element, the array keeps its length",
		Severity: SeverityWarning,
		Description: "In Bash, `unset arr[N]` removes the N-th element of the array (leaving a " +
			"sparse hole). In Zsh the same invocation passes the literal string `arr[N]` " +
			"to the `unset` builtin, which looks for a parameter with that name — finds " +
			"nothing — and returns success. The array is left untouched, `${#arr[@]}` " +
			"does not budge, and every downstream `for x in \"${arr[@]}\"` keeps iterating " +
			"the element the script thought it had removed. Use Zsh's native assignment " +
			"form `arr[N]=()` to delete an index, or `arr=(\"${(@)arr:#pattern}\")` to " +
			"filter by value.",
		Check: checkZC1856,
	})
}

func checkZC1856(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "unset" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1856IsArraySubscript(v) {
			return []Violation{{
				KataID: "ZC1856",
				Message: "`unset " + v + "` is a Bash idiom — in Zsh it tries to " +
					"unset a parameter literally named `" + v + "` and leaves the " +
					"array untouched. Use `arr[N]=()` or rebuild with " +
					"`arr=(\"${(@)arr:#pattern}\")`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1856IsArraySubscript(v string) bool {
	// Strip surrounding quotes and parser-applied `(...)` wrapping.
	trimmed := strings.TrimSpace(v)
	trimmed = strings.Trim(trimmed, "\"'")
	if strings.HasPrefix(trimmed, "(") && strings.HasSuffix(trimmed, ")") && len(trimmed) >= 2 {
		trimmed = trimmed[1 : len(trimmed)-1]
	}
	open := strings.Index(trimmed, "[")
	close := strings.LastIndex(trimmed, "]")
	if open <= 0 || close <= open+1 {
		return false
	}
	// The name portion must look like a shell identifier.
	return zc1856IsIdentifier(trimmed[:open])
}

func zc1856IsIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 && !zc1856IsIdentStart(r) {
			return false
		}
		if i > 0 && !zc1856IsIdentTail(r) {
			return false
		}
	}
	return true
}

func zc1856IsIdentStart(r rune) bool {
	return r == '_' || isAsciiLetter(r)
}

func zc1856IsIdentTail(r rune) bool {
	return r == '_' || isAsciiLetter(r) || (r >= '0' && r <= '9')
}

func isAsciiLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
