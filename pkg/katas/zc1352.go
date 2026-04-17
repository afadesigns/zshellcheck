package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1352",
		Title:    "Avoid `xargs -I{}` — use a Zsh `for` loop for per-item substitution",
		Severity: SeverityStyle,
		Description: "`xargs -I{}` runs one command per item with `{}` substituted. A Zsh `for` " +
			"loop over the same input (`for x in ${(f)\"$(cmd)\"}`) is clearer and keeps state " +
			"in the current shell.",
		Check: checkZC1352,
	})
}

func checkZC1352(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-I" || v == "--replace" {
			return []Violation{{
				KataID: "ZC1352",
				Message: "Avoid `xargs -I{}` — iterate with `for x in ${(f)\"$(cmd)\"}; do ...; done` " +
					"in Zsh. A for loop avoids one-subprocess-per-item and keeps variables in scope.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
