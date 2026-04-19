package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1758",
		Title:    "Warn on `gh codespace delete --force` — destroys codespace with uncommitted work",
		Severity: SeverityWarning,
		Description: "`gh codespace delete --force` (alias `-f`) skips the confirmation prompt " +
			"and deletes the target codespace along with any uncommitted, unpushed, or " +
			"unstaged work inside it. Combined with `--all`, one line wipes every codespace " +
			"on the account. Drop the flag, let the prompt enumerate what is about to go, " +
			"and only confirm after verifying no local state would be lost — `git status` " +
			"/ `git stash list` inside the codespace first.",
		Check: checkZC1758,
	})
}

func checkZC1758(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "codespace" || cmd.Arguments[1].String() != "delete" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if v == "--force" || v == "-f" {
			return []Violation{{
				KataID: "ZC1758",
				Message: "`gh codespace delete " + v + "` skips the prompt and drops " +
					"uncommitted work along with the codespace. Let the prompt list " +
					"what's about to go and verify `git status` inside first.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
