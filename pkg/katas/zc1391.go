// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1391",
		Title:    "Avoid `[[ -v VAR ]]` for Bash set-check — use Zsh `(( ${+VAR} ))`",
		Severity: SeverityWarning,
		Description: "Bash 4.2+ supports `[[ -v VAR ]]` to test whether a variable is set. Zsh " +
			"`[[ -v VAR ]]` is parsed but not as the set-check — Zsh's canonical form is " +
			"`(( ${+VAR} ))` which evaluates to 1 when set and 0 when unset, working reliably " +
			"across Zsh versions.",
		Check: checkZC1391,
	})
}

func checkZC1391(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	// `[[` is its own node type in most AST designs, but bracket-test tokens
	// may come through as commands. We look for "-v" as a bare arg with a
	// following identifier in a context that looks like a bracket test.
	if ident.Value != "test" && ident.Value != "[" && ident.Value != "[[" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() == "-v" && i+1 < len(cmd.Arguments) {
			next := cmd.Arguments[i+1].String()
			// Only flag if the "-v" is followed by an identifier (not a value to compare)
			if len(next) > 0 && !strings.Contains(next, "=") &&
				!strings.ContainsAny(next, "<>!/") {
				return []Violation{{
					KataID: "ZC1391",
					Message: "Use `(( ${+VAR} ))` for Zsh set-check — `-v` is a Bash 4.2+ " +
						"extension, not reliably portable to Zsh.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}

	return nil
}
