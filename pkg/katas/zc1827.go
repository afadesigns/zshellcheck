package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1827",
		Title:    "Error on `npm unpublish` — breaks every downstream that pinned the version",
		Severity: SeverityError,
		Description: "`npm unpublish PKG@VERSION` removes a published version from the registry. " +
			"Every downstream that pinned to that version — directly or through a transitive " +
			"lockfile entry — fails to install on the next `npm ci` / CI run. This is the " +
			"exact mechanism behind the 2016 `left-pad` outage; npm has since limited " +
			"unpublish to within 72 hours and added the `--force` gate, but within the " +
			"window the blast radius is still the whole ecosystem that pulled the package. " +
			"Use `npm deprecate PKG@VERSION 'reason'` instead — the version stays resolvable, " +
			"but installs print a warning and users can pin forward on their own schedule.",
		Check: checkZC1827,
	})
}

func checkZC1827(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "npm" && ident.Value != "pnpm" && ident.Value != "yarn" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "unpublish" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1827",
		Message: "`" + ident.Value + " unpublish` removes a published version — every " +
			"downstream that pinned it fails to install on next CI run (the left-pad " +
			"pattern). Use `npm deprecate PKG@VERSION 'reason'` so the version stays " +
			"resolvable with a warning.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
