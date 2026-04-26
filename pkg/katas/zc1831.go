// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1831SshUnits = map[string]bool{
	"ssh":            true,
	"sshd":           true,
	"ssh.service":    true,
	"sshd.service":   true,
	"ssh.socket":     true,
	"sshd.socket":    true,
	"openssh-server": true,
	"openssh":        true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1831",
		Title:    "Error on `systemctl stop|disable|mask ssh/sshd` — locks out the next remote login",
		Severity: SeverityError,
		Description: "Stopping, disabling, or masking the SSH daemon closes the door on the next " +
			"remote login. Existing connections survive for a while because sshd's spawned " +
			"per-session process keeps running, but any reconnect / CI follow-up step that " +
			"needs to ssh back in gets `Connection refused`. `systemctl disable ssh` and " +
			"`systemctl mask ssh` also survive reboots. Recovery requires console or out-of-" +
			"band access. If the goal is config reload, use `systemctl reload sshd`; if the " +
			"host is being retired, make sshd the last service you touch.",
		Check: checkZC1831,
	})
}

func checkZC1831(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	action := cmd.Arguments[0].String()
	if action != "stop" && action != "disable" && action != "mask" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		unit := arg.String()
		if zc1831SshUnits[unit] {
			return []Violation{{
				KataID: "ZC1831",
				Message: "`systemctl " + action + " " + unit + "` blocks SSH — " +
					"existing sessions survive but reconnects fail. `disable`/`mask` " +
					"persist across reboots. Use `reload sshd` for config changes.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
