// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1629",
		Title:    "Warn on `rsync --rsync-path='sudo rsync'` — hidden remote privilege escalation",
		Severity: SeverityWarning,
		Description: "`--rsync-path` normally overrides the path to the remote rsync binary. " +
			"Setting it to `sudo rsync` (or `doas rsync` / `pkexec rsync`) instead makes the " +
			"remote side run rsync as root. That is sometimes legitimate — copying into " +
			"`/etc/` from a CI job — but the flag is easy to miss in review because it looks " +
			"like a path override. Provision a scoped sudoers rule that names exactly which " +
			"rsync invocation the remote user may run, and keep the path explicit (`--rsync-" +
			"path=/usr/bin/rsync`).",
		Check: checkZC1629,
	})
}

func checkZC1629(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "rsync" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !strings.HasPrefix(v, "--rsync-path=") {
			continue
		}
		val := strings.TrimPrefix(v, "--rsync-path=")
		val = strings.Trim(val, "\"'")
		if strings.Contains(val, "sudo") ||
			strings.Contains(val, "doas") ||
			strings.Contains(val, "pkexec") {
			return []Violation{{
				KataID: "ZC1629",
				Message: "`rsync --rsync-path='" + val + "'` runs remote rsync under " +
					"privilege escalation. Use a scoped sudoers rule on the remote host " +
					"and keep the path explicit (`/usr/bin/rsync`).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
