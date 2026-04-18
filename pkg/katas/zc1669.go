package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1669",
		Title:    "Warn on `git gc --prune=now` / `git reflog expire --expire=now` — deletes recovery window",
		Severity: SeverityWarning,
		Description: "Git keeps dropped commits and orphaned objects for `gc.reflogExpire` " +
			"(default 90 days) and `gc.pruneExpire` (default two weeks) so a `git reflog` + " +
			"`git reset` can still recover work you thought you threw away. `git gc " +
			"--prune=now` and `git reflog expire --expire=now --all` bulldoze both windows " +
			"in one go — a stray interactive rebase no longer has a safety net. Use the " +
			"default cadence (`git gc`, no `--prune=now`) unless you are actively purging " +
			"leaked secrets or proof-of-concept code; pair the destructive form with a " +
			"stale mirror push so at least one copy of the dropped history remains.",
		Check: checkZC1669,
	})
}

func checkZC1669(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	switch cmd.Arguments[0].String() {
	case "gc":
		for _, arg := range cmd.Arguments[1:] {
			if arg.String() == "--prune=now" || arg.String() == "--prune=0" {
				return zc1669Hit(cmd, "git gc --prune=now")
			}
		}
	case "reflog":
		if len(cmd.Arguments) < 2 || cmd.Arguments[1].String() != "expire" {
			return nil
		}
		for _, arg := range cmd.Arguments[2:] {
			if arg.String() == "--expire=now" {
				return zc1669Hit(cmd, "git reflog expire --expire=now")
			}
		}
	}
	return nil
}

func zc1669Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1669",
		Message: "`" + form + "` bulldozes the reflog / prune recovery window — keep the " +
			"default cadence unless you are actively purging leaked secrets, and mirror " +
			"the dropped history off-box first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
