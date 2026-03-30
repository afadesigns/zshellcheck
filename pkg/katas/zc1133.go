package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1133",
		Title: "Avoid `kill -9` — use `kill` first, then escalate",
		Description: "`kill -9` (SIGKILL) cannot be caught or ignored. Always try `kill` (SIGTERM) first " +
			"to allow the process to clean up, then use `kill -9` only as a last resort.",
		Severity: SeverityStyle,
		Check:    checkZC1133,
	})
}

func checkZC1133(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kill" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-9" || val == "-KILL" || val == "-SIGKILL" {
			return []Violation{{
				KataID: "ZC1133",
				Message: "Avoid `kill -9` as a first resort. Use `kill` (SIGTERM) first " +
					"to allow graceful shutdown, then escalate to `kill -9` if needed.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
