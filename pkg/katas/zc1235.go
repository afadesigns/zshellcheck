package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1235",
		Title:    "Use `git push --force-with-lease` instead of `--force`",
		Severity: SeverityWarning,
		Description: "`git push --force` overwrites remote history unconditionally. " +
			"`--force-with-lease` is safer as it fails if the remote has changed.",
		Check: checkZC1235,
	})
}

func checkZC1235(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "push" {
		return nil
	}

	hasForce := false
	hasFWL := false

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "-f" {
			hasForce = true
		}
	}

	if hasForce && !hasFWL {
		return []Violation{{
			KataID: "ZC1235",
			Message: "Use `git push --force-with-lease` instead of `-f`/`--force`. " +
				"It prevents overwriting remote changes made by others.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
