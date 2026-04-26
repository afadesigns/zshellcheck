// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1149",
		Title:    "Avoid `echo` for error messages — use `>&2`",
		Severity: SeverityInfo,
		Description: "Error messages should go to stderr, not stdout. " +
			"Use `print -u2` or `echo ... >&2` to ensure errors are properly separated.",
		Check: checkZC1149,
	})
}

func checkZC1149(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "echo" && ident.Value != "printf" && ident.Value != "print" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		// Check for error-like messages
		if len(val) > 5 {
			clean := val
			if len(clean) > 2 && (clean[0] == '\'' || clean[0] == '"') {
				clean = clean[1 : len(clean)-1]
			}
			if len(clean) >= 5 && (clean[:5] == "Error" || clean[:5] == "error" || clean[:5] == "ERROR") {
				return []Violation{{
					KataID: "ZC1149",
					Message: "Error messages should go to stderr. Use `print -u2` or append `>&2` " +
						"to separate error output from normal stdout.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityInfo,
				}}
			}
		}
	}

	return nil
}
