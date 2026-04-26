// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1511",
		Title:    "Error on `nmcli ... <wireless/vpn secret>` on command line",
		Severity: SeverityError,
		Description: "Passing Wi-Fi pre-shared keys or VPN secrets as positional `nmcli` args " +
			"puts them in `ps`, shell history, and `/proc/<pid>/cmdline`. Let NetworkManager " +
			"store the secret for you via `--ask` (interactive prompt, no TTY echo) or use " +
			"`keyfile` connection profiles under `/etc/NetworkManager/system-connections/` " +
			"with mode 0600.",
		Check: checkZC1511,
	})
}

var nmcliSecretKeys = []string{
	"802-11-wireless-security.psk",
	"wifi-sec.psk",
	"wifi.psk",
	"vpn.secrets.password",
	"ipsec-secret",
	"openvpn-password",
	"802-1x.password",
}

func checkZC1511(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "nmcli" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	for i, a := range args {
		low := strings.ToLower(a)
		for _, key := range nmcliSecretKeys {
			if low == key && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				return []Violation{{
					KataID: "ZC1511",
					Message: "`nmcli` passed `" + key + " <secret>` on the command line — " +
						"ends up in ps/history. Use `--ask` or a keyfile profile.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
	}
	return nil
}
