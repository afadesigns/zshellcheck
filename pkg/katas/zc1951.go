package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1951",
		Title:    "Error on `ceph osd pool delete … --yes-i-really-really-mean-it` — automates Ceph's double-safety phrase",
		Severity: SeverityError,
		Description: "Ceph intentionally requires both the pool name twice and the flag " +
			"`--yes-i-really-really-mean-it` before it will delete a pool, so a typo during a " +
			"live operation cannot drop production data. Baking the phrase into a script " +
			"defeats the friction — a rebase of the wrong variable, a typo in the pool name, " +
			"or a stale `for pool in $(…)` loop then silently deletes real pools. Remove the " +
			"flag from scripts. Do the deletion interactively, or wrap it in a runbook that " +
			"spells out the pool name in the commit message the operator acknowledges.",
		Check: checkZC1951,
	})
}

func checkZC1951(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `ceph osd pool delete … --yes-i-really-really-mean-it`
	// mangles the command name to `yes-i-really-really-mean-it`.
	if ident.Value == "yes-i-really-really-mean-it" ||
		ident.Value == "yes-i-really-mean-it" {
		return zc1951Hit(cmd)
	}
	if ident.Value != "ceph" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--yes-i-really-really-mean-it" || v == "--yes-i-really-mean-it" {
			return zc1951Hit(cmd)
		}
	}
	return nil
}

func zc1951Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1951",
		Message: "`ceph … --yes-i-really-really-mean-it` automates the double-safety " +
			"phrase — a typo or stale loop silently deletes production pools. Run " +
			"deletions interactively, or spell the pool name in a runbook commit.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
