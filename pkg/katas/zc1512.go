// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1512",
		Title:    "Style: `service <unit> <verb>` — use `systemctl <verb> <unit>` on systemd hosts",
		Severity: SeverityStyle,
		Description: "`service` is the SysV init compatibility wrapper. On a systemd-managed " +
			"host (every mainstream distro since ~2016) it translates to `systemctl` anyway, " +
			"but reverses argument order, loses `--user` scope, ignores unit templating, and " +
			"can't restart sockets or timers. Prefer `systemctl start|stop|restart|reload " +
			"<unit>` for consistency across scripts and interactive shells.",
		Check: checkZC1512,
		Fix:   fixZC1512,
	})
}

// fixZC1512 rewrites `service UNIT VERB` into `systemctl VERB UNIT`.
// Three edits per match: rename `service` → `systemctl`, swap the
// textual contents of the UNIT and VERB positions. Gated to simple
// Identifier args so the swap stays byte-exact; concat-form units
// (rare in practice) stay detection-only.
func fixZC1512(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "service" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	unitIdent, ok := cmd.Arguments[0].(*ast.Identifier)
	if !ok {
		return nil
	}
	verbIdent, ok := cmd.Arguments[1].(*ast.Identifier)
	if !ok {
		return nil
	}
	unitTok := unitIdent.TokenLiteralNode()
	verbTok := verbIdent.TokenLiteralNode()
	unitOff := LineColToByteOffset(source, unitTok.Line, unitTok.Column)
	verbOff := LineColToByteOffset(source, verbTok.Line, verbTok.Column)
	if unitOff < 0 || verbOff < 0 {
		return nil
	}
	if unitOff+len(unitIdent.Value) > len(source) ||
		string(source[unitOff:unitOff+len(unitIdent.Value)]) != unitIdent.Value {
		return nil
	}
	if verbOff+len(verbIdent.Value) > len(source) ||
		string(source[verbOff:verbOff+len(verbIdent.Value)]) != verbIdent.Value {
		return nil
	}
	return []FixEdit{
		{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("service"),
			Replace: "systemctl",
		},
		{
			Line:    unitTok.Line,
			Column:  unitTok.Column,
			Length:  len(unitIdent.Value),
			Replace: verbIdent.Value,
		},
		{
			Line:    verbTok.Line,
			Column:  verbTok.Column,
			Length:  len(verbIdent.Value),
			Replace: unitIdent.Value,
		},
	}
}

func checkZC1512(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "service" {
		return nil
	}

	// Needs at least <unit> <verb>.
	if len(cmd.Arguments) < 2 {
		return nil
	}
	verb := cmd.Arguments[1].String()
	switch verb {
	case "start", "stop", "restart", "reload", "status", "force-reload", "try-restart":
	default:
		return nil
	}

	unit := cmd.Arguments[0].String()
	return []Violation{{
		KataID: "ZC1512",
		Message: "`service " + unit + " " + verb + "` — prefer `systemctl " + verb + " " +
			unit + "` for consistency with other systemd commands.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
