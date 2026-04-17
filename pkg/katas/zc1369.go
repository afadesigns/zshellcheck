package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1369",
		Title:    "Prefer Zsh `${(V)var}` over `od -c` for printable-visible character output",
		Severity: SeverityStyle,
		Description: "Zsh's `${(V)var}` parameter flag renders non-printable characters in " +
			"visible form (e.g. `\\n` for newline). For simple inspection of a variable's " +
			"contents, this avoids the `od -c` process entirely.",
		Check: checkZC1369,
	})
}

func checkZC1369(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "od" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" || v == "-C" {
			return []Violation{{
				KataID: "ZC1369",
				Message: "Use Zsh `${(V)var}` to see non-printable characters in a variable — " +
					"renders control chars as `\\n`, `\\t`, etc., without spawning `od`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
