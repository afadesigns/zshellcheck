// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1303",
		Title:    "Avoid `enable` command — use `zmodload` for Zsh modules",
		Severity: SeverityWarning,
		Description: "The `enable` command is a Bash builtin for enabling/disabling builtins. " +
			"Zsh uses `zmodload` to load and manage modules, and `disable`/`enable` " +
			"have different semantics. Use `zmodload` for module management.",
		Check: checkZC1303,
	})
}

func checkZC1303(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "enable" {
		return nil
	}

	// enable with -f flag loads a builtin from a shared object (Bash-specific)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-f" {
			return []Violation{{
				KataID:  "ZC1303",
				Message: "Avoid `enable -f` in Zsh — use `zmodload` to load modules. `enable -f` is Bash-specific.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	}

	return nil
}
