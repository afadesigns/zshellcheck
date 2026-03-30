package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1135",
		Title: "Avoid `env VAR=val cmd` — use inline assignment",
		Description: "Zsh supports inline environment variable assignment with `VAR=val cmd`. " +
			"Avoid spawning `env` for simple variable-prefixed command execution.",
		Severity: SeverityStyle,
		Check:    checkZC1135,
	})
}

func checkZC1135(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "env" {
		return nil
	}

	// Only flag env with VAR=val patterns followed by a command
	// Skip env -i (clean environment), env -u (unset), env -S
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	// Check if any argument contains = (env var assignment)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if strings.Contains(val, "=") {
			return []Violation{{
				KataID: "ZC1135",
				Message: "Use inline `VAR=val cmd` instead of `env VAR=val cmd`. " +
					"Zsh supports inline env assignment without spawning env.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
