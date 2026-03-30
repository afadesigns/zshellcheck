package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1172",
		Title:    "Use `read -A` instead of Bash `read -a` for arrays",
		Severity: SeverityInfo,
		Description: "Bash uses `read -a` to read into an array, but Zsh uses `read -A`. " +
			"Using `-a` in Zsh reads into a scalar, not an array.",
		Check: checkZC1172,
	})
}

func checkZC1172(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "read" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-a" {
			return []Violation{{
				KataID: "ZC1172",
				Message: "Use `read -A` instead of `read -a` in Zsh. " +
					"The `-a` flag is Bash syntax; Zsh uses `-A` to read into arrays.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}
