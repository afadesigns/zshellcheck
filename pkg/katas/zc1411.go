package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1411",
		Title:    "Use Zsh `disable` instead of Bash `enable -n` to hide builtins",
		Severity: SeverityStyle,
		Description: "Bash's `enable -n name` disables a builtin so that the external of the same " +
			"name is used. Zsh provides a dedicated `disable` builtin: `disable name` achieves " +
			"the same in one verb. Re-enable later with `enable name`.",
		Check: checkZC1411,
	})
}

func checkZC1411(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "enable" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-n" {
			return []Violation{{
				KataID: "ZC1411",
				Message: "Use Zsh `disable name` instead of `enable -n name`. Zsh has a " +
					"dedicated `disable` builtin that reads clearer.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
