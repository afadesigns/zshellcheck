package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1183",
		Title:    "Use Zsh glob qualifiers instead of `ls -t` for file ordering",
		Severity: SeverityStyle,
		Description: "Zsh glob qualifiers like `*(om[1])` (newest) or `*(Om[1])` (oldest) " +
			"order files without spawning `ls`. Avoid `ls -t | head` patterns.",
		Check: checkZC1183,
	})
}

func checkZC1183(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ls" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-t" || val == "-tr" || val == "-lt" || val == "-ltr" {
			return []Violation{{
				KataID: "ZC1183",
				Message: "Use Zsh glob qualifiers `*(om[1])` for newest file or `*(Om[1])` for oldest " +
					"instead of `ls -t`. Glob qualifiers avoid spawning external processes.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
