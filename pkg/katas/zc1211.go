package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1211",
		Title:    "Use `git stash push -m` instead of bare `git stash`",
		Severity: SeverityStyle,
		Description: "Bare `git stash` creates unnamed stashes that are hard to identify later. " +
			"Use `git stash push -m 'description'` for self-documenting stashes.",
		Check: checkZC1211,
	})
}

func checkZC1211(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 {
		return nil
	}

	subCmd := cmd.Arguments[0].String()
	if subCmd != "stash" {
		return nil
	}

	// git stash push -m is fine, git stash pop/apply/list/drop are fine
	if len(cmd.Arguments) >= 2 {
		action := cmd.Arguments[1].String()
		if action == "push" || action == "pop" || action == "apply" ||
			action == "list" || action == "drop" || action == "show" {
			return nil
		}
	}

	// Bare "git stash" with no subcommand
	if len(cmd.Arguments) == 1 {
		return []Violation{{
			KataID: "ZC1211",
			Message: "Use `git stash push -m 'description'` instead of bare `git stash`. " +
				"Named stashes are easier to identify and manage.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
