// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1613Readers = map[string]bool{
	"cat": true, "less": true, "more": true,
	"head": true, "tail": true, "wc": true,
	"grep": true, "awk": true, "sed": true, "cut": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1613",
		Title:    "Warn on reading SSH private-key files with `cat` / `less` / `grep` / `head`",
		Severity: SeverityWarning,
		Description: "Piping an SSH private key through a generic text tool copies the raw " +
			"key material into the process and — if stdout is redirected or piped — often " +
			"into logs, backup files, or a terminal scrollback buffer. Host keys under " +
			"`/etc/ssh/ssh_host_*_key` impersonate the server; user keys under `~/.ssh/id_*` " +
			"impersonate the user. Use `ssh-keygen -l -f KEY` for fingerprint / metadata, or " +
			"pass the key path to the consumer directly (`ssh -i`, `git -c core.sshCommand`) " +
			"without staging it through a shell tool.",
		Check: checkZC1613,
	})
}

func checkZC1613(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if !zc1613Readers[ident.Value] {
		return nil
	}

	userSuffixes := []string{
		"/.ssh/id_rsa", "/.ssh/id_ed25519", "/.ssh/id_ecdsa", "/.ssh/id_dsa",
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "/etc/ssh/ssh_host_") && strings.HasSuffix(v, "_key") {
			return zc1613Hit(cmd, v)
		}
		for _, s := range userSuffixes {
			if strings.HasSuffix(v, s) {
				return zc1613Hit(cmd, v)
			}
		}
	}
	return nil
}

func zc1613Hit(cmd *ast.SimpleCommand, path string) []Violation {
	return []Violation{{
		KataID: "ZC1613",
		Message: "Reading `" + path + "` through a text tool copies private-key material " +
			"into the process and often into logs / scrollback. Use `ssh-keygen -l -f` " +
			"for metadata, or pass the path directly to the consumer.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
