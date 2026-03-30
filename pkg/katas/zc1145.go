package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1145",
		Title:    "Avoid `tr -d` for character deletion — use parameter expansion",
		Severity: SeverityStyle,
		Description: "For simple character deletion from variables, use Zsh `${var//char/}` " +
			"instead of piping through `tr -d`.",
		Check: checkZC1145,
	})
}

func checkZC1145(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tr" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	firstArg := cmd.Arguments[0].String()
	if firstArg == "-d" {
		// Check second arg is a simple single char
		secondArg := cmd.Arguments[1].String()
		if len(secondArg) <= 3 { // Simple char like 'x' or " "
			return []Violation{{
				KataID: "ZC1145",
				Message: "Use `${var//char/}` instead of piping through `tr -d`. " +
					"Parameter expansion is faster for simple character deletion.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
