package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1178",
		Title:    "Avoid `stty` for terminal size — use Zsh `$COLUMNS`/`$LINES`",
		Severity: SeverityStyle,
		Description: "Zsh maintains `$COLUMNS` and `$LINES` as built-in variables tracking " +
			"terminal dimensions. Avoid spawning `stty` or `tput` for size queries.",
		Check: checkZC1178,
	})
}

func checkZC1178(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "stty" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "size" {
			return []Violation{{
				KataID: "ZC1178",
				Message: "Use `$COLUMNS` and `$LINES` instead of `stty size`. " +
					"Zsh tracks terminal dimensions as built-in variables.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
