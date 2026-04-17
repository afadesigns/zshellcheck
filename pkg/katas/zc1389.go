package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1389",
		Title:    "Avoid `$HOSTFILE` — Bash-only; Zsh uses `$hosts` array",
		Severity: SeverityWarning,
		Description: "Bash reads `$HOSTFILE` to feed hostname completion. Zsh populates hostname " +
			"completion from the `$hosts` array (lowercase). Setting `$HOSTFILE` in Zsh is " +
			"ignored; extend `$hosts` instead.",
		Check: checkZC1389,
	})
}

func checkZC1389(node ast.Node) []Violation {
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
		if strings.Contains(v, "HOSTFILE") {
			return []Violation{{
				KataID: "ZC1389",
				Message: "`$HOSTFILE` is Bash-only. Zsh reads hostnames for completion from the " +
					"`$hosts` array (lowercase).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
