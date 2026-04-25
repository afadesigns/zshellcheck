package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1217",
		Title:    "Avoid `service` command — use `systemctl` on systemd",
		Severity: SeverityInfo,
		Description: "`service` is a SysVinit compatibility wrapper. " +
			"On systemd systems, use `systemctl start/stop/restart/status` directly.",
		Check: checkZC1217,
		// Reuse the `service UNIT VERB` → `systemctl VERB UNIT` rewrite
		// from ZC1512. Both detectors fire on the same shape; the
		// conflict resolver dedupes overlapping edits.
		Fix: fixZC1512,
	})
}

func checkZC1217(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "service" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1217",
		Message: "Avoid `service` — it is a SysVinit compatibility wrapper. " +
			"Use `systemctl` directly on systemd systems.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}
