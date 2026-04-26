// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1680",
		Title:    "Error on `ansible-playbook --vault-password-file=/tmp/...` — world-traversable vault key",
		Severity: SeverityError,
		Description: "The Ansible Vault decryption key lives in the `--vault-password-file` " +
			"path. `/tmp`, `/var/tmp`, and `/dev/shm` are world-traversable: a concurrent " +
			"local user who guesses (or `inotifywait`s for) the filename opens it during " +
			"the playbook run and dumps every secret the vault protects. Keep vault keys " +
			"in a root-owned mode-0400 file under `/etc/ansible/` or `$HOME/.ansible/`, or " +
			"supply the passphrase via a no-echo helper script (`vault-password-client`) " +
			"that fetches from `pass` / `vault kv get`.",
		Check: checkZC1680,
	})
}

func checkZC1680(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ansible-playbook" && ident.Value != "ansible" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "--vault-password-file=") {
			if zc1680Unsafe(strings.TrimPrefix(v, "--vault-password-file=")) {
				return zc1680Hit(cmd)
			}
			continue
		}
		if v == "--vault-password-file" && i+1 < len(cmd.Arguments) {
			if zc1680Unsafe(cmd.Arguments[i+1].String()) {
				return zc1680Hit(cmd)
			}
		}
	}
	return nil
}

func zc1680Unsafe(path string) bool {
	return strings.HasPrefix(path, "/tmp/") ||
		strings.HasPrefix(path, "/var/tmp/") ||
		strings.HasPrefix(path, "/dev/shm/")
}

func zc1680Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1680",
		Message: "`ansible-playbook --vault-password-file` under `/tmp/` / `/var/tmp/` / " +
			"`/dev/shm/` — world-traversable, any local user can race-read it. Store the " +
			"key mode-0400 under `/etc/ansible/` or supply via a `vault-password-client` helper.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
