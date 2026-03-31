package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1191",
		Title:    "Use `$+commands[cmd]` for command existence in Zsh",
		Severity: SeverityStyle,
		Description: "Zsh maintains `$commands` associative array of all available commands. " +
			"Use `(( $+commands[cmd] ))` instead of `command -v cmd > /dev/null` patterns.",
		Check: checkZC1191,
	})
}

func checkZC1191(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "command" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	first := cmd.Arguments[0].String()
	if first == "-v" {
		return []Violation{{
			KataID: "ZC1191",
			Message: "Use `(( $+commands[cmd] ))` instead of `command -v cmd`. " +
				"Zsh `$commands` array provides instant command lookups.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
