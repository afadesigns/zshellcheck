package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1765",
		Title:    "Error on `snap remove --purge SNAP` — skips the automatic data snapshot",
		Severity: SeverityError,
		Description: "`snap remove SNAP` takes a snapshot of every writable area (`$SNAP_DATA`, " +
			"`$SNAP_USER_DATA`, `$SNAP_COMMON`) before uninstalling, so the data can later " +
			"be restored with `snap restore`. `--purge` skips that snapshot: the snap is " +
			"gone along with every file it owned, and snapd has no record to roll back. " +
			"Drop `--purge` unless the snap's data is genuinely disposable; otherwise " +
			"`snap save SNAP` first, capture the set ID, and only then remove.",
		Check: checkZC1765,
	})
}

func checkZC1765(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "snap" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "remove" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "--purge" {
			return []Violation{{
				KataID: "ZC1765",
				Message: "`snap remove --purge` skips the pre-remove data snapshot — the " +
					"snap's files are gone with no rollback. Drop `--purge` or capture " +
					"a `snap save` set ID first.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
