package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1189",
		Title:    "Avoid `source /dev/stdin` — use direct evaluation",
		Severity: SeverityWarning,
		Description: "`source /dev/stdin` is fragile and platform-dependent. " +
			"Use `eval \"$(cmd)\"` or direct command execution instead.",
		Check: checkZC1189,
	})
}

func checkZC1189(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "source" && ident.Value != "." {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/dev/stdin" || val == "/proc/self/fd/0" {
			return []Violation{{
				KataID: "ZC1189",
				Message: "Avoid `source /dev/stdin`. Use `eval \"$(cmd)\"` for direct evaluation. " +
					"`/dev/stdin` sourcing is fragile across platforms.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
