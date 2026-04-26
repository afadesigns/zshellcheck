// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1500",
		Title:    "Warn on `systemctl edit <unit>` in scripts — requires interactive editor",
		Severity: SeverityWarning,
		Description: "`systemctl edit <unit>` (without `--no-edit` and without a piped `EDITOR`) " +
			"opens `$EDITOR` on a tmpfile and waits for the user. In a non-interactive script " +
			"this either hangs until timeout or silently succeeds with no change, depending on " +
			"how the editor handles a closed stdin. For scripted unit tweaks, drop a `.conf` " +
			"drop-in under `/etc/systemd/system/<unit>.d/` and call `systemctl daemon-reload`.",
		Check: checkZC1500,
	})
}

func checkZC1500(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "systemctl" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "edit" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--no-edit" || v == "--runtime" {
			// Still odd, but at least doesn't spin on an editor. Let it pass.
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1500",
		Message: "`systemctl edit` opens $EDITOR and waits for the user. Use a drop-in " +
			"`/etc/systemd/system/<unit>.d/*.conf` + `daemon-reload` in scripts.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
