package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1360",
		Title:    "Use Zsh `*(OL)` glob qualifier instead of `ls -S` for size-ordered listing",
		Severity: SeverityStyle,
		Description: "Zsh glob qualifier `*(OL)` orders results by size (descending). `*(oL)` is " +
			"ascending. Combined with `[N]` subscript you get the N-th largest/smallest file " +
			"without `ls -S` and piping.",
		Check: checkZC1360,
	})
}

func checkZC1360(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ls" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-S" || v == "-Sr" || v == "-lS" || v == "-lSr" {
			return []Violation{{
				KataID: "ZC1360",
				Message: "Use Zsh `*(OL)` (largest-first) or `*(oL)` (smallest-first) glob qualifier " +
					"instead of `ls -S`. No external process needed.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
