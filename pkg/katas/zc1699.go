package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1699",
		Title:    "Warn on `kubectl drain --delete-emptydir-data` — pod-local scratch data lost",
		Severity: SeverityWarning,
		Description: "`kubectl drain NODE --delete-emptydir-data` (older alias `--delete-local-" +
			"data`) lets drain evict pods that mount an `emptyDir` volume — the volume is " +
			"deleted along with the pod, destroying any scratch data it held. Production " +
			"clusters use `emptyDir` for caches, write-ahead logs, and scratch state that " +
			"takes hours to rebuild. Confirm the pods on the node tolerate the loss (or " +
			"migrate to a `persistentVolumeClaim`) before adding the flag; otherwise plan " +
			"a controlled drain without it and accept the stuck-drain warning for the " +
			"affected pods.",
		Check: checkZC1699,
	})
}

func checkZC1699(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "drain" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--delete-emptydir-data" || v == "--delete-local-data" {
			return []Violation{{
				KataID: "ZC1699",
				Message: "`kubectl drain " + v + "` deletes `emptyDir` volumes along with the " +
					"evicted pods — caches / WAL / scratch state are lost. Verify tolerance " +
					"or migrate to a PersistentVolumeClaim first.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
