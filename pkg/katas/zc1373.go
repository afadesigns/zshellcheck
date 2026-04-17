package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1373",
		Title:    "Use Zsh `${(0)var}` flag for NUL-split parsing instead of `env -0`",
		Severity: SeverityStyle,
		Description: "When reading NUL-terminated data (e.g. `/proc/*/environ`), Zsh's `${(0)var}` " +
			"parameter flag splits on NUL into an array natively. Avoid `env -0 | xargs -0 ...` " +
			"chains that require two additional processes.",
		Check: checkZC1373,
	})
}

func checkZC1373(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "env" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-0" || v == "--null" {
			return []Violation{{
				KataID: "ZC1373",
				Message: "Use Zsh `${(0)\"$(<file)\"}` to split NUL-terminated content in-shell. " +
					"`env -0` is usually followed by `xargs -0` or a read loop — both avoided.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
