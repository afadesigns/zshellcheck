package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1236",
		Title:    "Avoid `git reset --hard` — irreversible data loss risk",
		Severity: SeverityWarning,
		Description: "`git reset --hard` discards all uncommitted changes irreversibly. " +
			"Use `git stash` to save changes first, or `git reset --soft` to keep them staged.",
		Check: checkZC1236,
	})
}

func checkZC1236(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "reset" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--hard" {
			return []Violation{{
				KataID: "ZC1236",
				Message: "Avoid `git reset --hard` — it permanently discards uncommitted changes. " +
					"Use `git stash` first, or `git reset --soft` to keep changes staged.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
