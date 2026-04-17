package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1433",
		Title:    "Caution with `userdel -f` / `-r` — removes home directory and kills processes",
		Severity: SeverityWarning,
		Description: "`userdel -f` proceeds even when the user is logged in or has running " +
			"processes, potentially killing unsaved work. `-r` additionally deletes the home " +
			"directory and mail spool. Combined (`-rf`) these are destructive and often " +
			"misused for 'clean up a user' without warning. Verify no active sessions first.",
		Check: checkZC1433,
	})
}

func checkZC1433(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "userdel" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" || v == "-r" || v == "-rf" || v == "-fr" ||
			v == "--force" || v == "--remove" {
			return []Violation{{
				KataID: "ZC1433",
				Message: "`userdel -f`/`-r` forcibly removes user (kills processes, deletes home). " +
					"Check for active sessions first with `who -u` / `loginctl list-users`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
