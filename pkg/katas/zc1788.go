// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1788SSHMutablePrefixes = []string{
	"/tmp/",
	"/var/tmp/",
	"/dev/shm/",
	"/home/",
	"/root/",
	"/opt/",
	"/srv/",
	"/mnt/",
	"/media/",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1788",
		Title:    "Warn on `ssh -F /tmp/config` — config from a mutable path can pin `ProxyCommand` to arbitrary code",
		Severity: SeverityWarning,
		Description: "`ssh -F PATH` (and `scp -F PATH`, `sftp -F PATH`) loads a user-supplied " +
			"config file. Anything in `/etc/ssh/ssh_config` can be overridden — notably " +
			"`ProxyCommand`, `LocalCommand`, `PermitLocalCommand`, and `Include` — which means " +
			"a mutable source path is an execution primitive: another local user flips " +
			"`ProxyCommand` to `/tmp/pwn`, and the next `ssh` run launches it with the " +
			"caller's credentials and forwarded agent. Keep the config in `~/.ssh/config` (or " +
			"a repo-owned path with the same owner and `0600` perms) and never pass `-F` to " +
			"`/tmp`, `/var/tmp`, `/dev/shm`, another user's `/home`, `/opt`, `/srv`, or `/mnt`.",
		Check: checkZC1788,
	})
}

func checkZC1788(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" && ident.Value != "scp" && ident.Value != "sftp" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		var path string
		switch {
		case v == "-F":
			if i+1 >= len(cmd.Arguments) {
				return nil
			}
			path = cmd.Arguments[i+1].String()
		case strings.HasPrefix(v, "-F"):
			path = v[2:]
		default:
			continue
		}
		path = strings.Trim(path, "\"'")
		if !strings.HasPrefix(path, "/") {
			continue
		}
		for _, prefix := range zc1788SSHMutablePrefixes {
			if strings.HasPrefix(path, prefix) {
				return []Violation{{
					KataID: "ZC1788",
					Message: "`" + ident.Value + " -F " + path + "` loads an alternate " +
						"config from a mutable path — a tamper on that file can pin " +
						"`ProxyCommand` to arbitrary code. Keep the config in " +
						"`~/.ssh/config` (or a repo-owned path with `0600` perms).",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}
