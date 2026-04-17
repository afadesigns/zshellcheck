package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1416",
		Title:    "Prefer Zsh `preexec` hook over `trap 'cmd' DEBUG`",
		Severity: SeverityWarning,
		Description: "Bash's `trap 'cmd' DEBUG` runs `cmd` before each simple command. Zsh's " +
			"equivalent is the `preexec` function (or `add-zsh-hook preexec name`) which " +
			"receives the about-to-execute command line as `$1`, `$2`, `$3`. The DEBUG trap " +
			"is not fired in Zsh the way it is in Bash — use preexec for portability.",
		Check: checkZC1416,
	})
}

func checkZC1416(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "trap" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "DEBUG" {
			return []Violation{{
				KataID: "ZC1416",
				Message: "Use Zsh `preexec() { ... }` (or `add-zsh-hook preexec`) instead of " +
					"`trap 'cmd' DEBUG`. Zsh's DEBUG trap does not fire the same way as Bash's.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
