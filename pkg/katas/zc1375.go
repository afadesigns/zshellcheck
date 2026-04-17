package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1375",
		Title:    "Use `set --` (empty) instead of `shift $#` to clear positional arguments",
		Severity: SeverityStyle,
		Description: "`shift $#` discards all positional arguments but is indirect and depends on " +
			"`$#` being accurate at evaluation. `set --` (with nothing after `--`) clears the " +
			"positional parameters atomically and reads more clearly.",
		Check: checkZC1375,
	})
}

func checkZC1375(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "shift" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "$#" || v == "${#}" {
			return []Violation{{
				KataID: "ZC1375",
				Message: "Use `set --` instead of `shift $#` to clear positional arguments. " +
					"Clearer intent, no dependency on `$#` accuracy at evaluation time.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
