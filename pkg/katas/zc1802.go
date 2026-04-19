package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1802",
		Title:    "Warn on `dnf history undo N` / `rollback N` — reverses transactions without compat check",
		Severity: SeverityWarning,
		Description: "`dnf history undo N` reverts the exact package set of transaction N — every " +
			"install turns into a remove, every remove into an install, every update into a " +
			"downgrade. `dnf history rollback N` does the same for every transaction after " +
			"N. Neither checks that the older versions still resolve cleanly against the " +
			"current package graph: dependencies that moved forward for other reasons end up " +
			"downgraded alongside, security patches get reverted, and services whose " +
			"configuration was migrated fail to start on the older binary. Review the plan " +
			"with `dnf history info N`, pin the rollback scope with `--exclude=` / `--assumeyes` " +
			"only after review, or restore from a filesystem snapshot taken before the " +
			"original transaction.",
		Check: checkZC1802,
	})
}

func checkZC1802(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dnf" && ident.Value != "yum" && ident.Value != "dnf5" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "history" {
		return nil
	}
	action := cmd.Arguments[1].String()
	if action != "undo" && action != "rollback" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1802",
		Message: "`" + ident.Value + " history " + action + "` reverses the past " +
			"transaction — deps downgrade, security patches can get reverted. " +
			"Review with `dnf history info`, or restore a filesystem snapshot.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
