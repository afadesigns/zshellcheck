// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1533",
		Title:    "Warn on `setsid <cmd>` — detaches from controlling TTY, escapes supervision",
		Severity: SeverityWarning,
		Description: "`setsid` starts a new session and process group. Combined with `-f` " +
			"(`--fork`) the child is fully detached from the invoking shell: `SIGHUP` from " +
			"logout does not reach it, the tty hang-up no longer terminates it, and it falls " +
			"off the script's job table. That is legitimate for daemonising a long-running " +
			"helper (though systemd does this better) and is also a standard persistence " +
			"mechanism. Prefer a systemd unit; if you must detach, document why.",
		Check: checkZC1533,
	})
}

func checkZC1533(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setsid" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1533",
		Message: "`setsid` detaches the child from the TTY / session — escapes supervision. " +
			"Prefer a systemd unit; document a detach if one is genuinely needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
