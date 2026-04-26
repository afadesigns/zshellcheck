// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1595",
		Title:    "Warn on `setfacl -m u:nobody:... / o::rwx` — ACL grants that bypass `chmod` scrutiny",
		Severity: SeverityWarning,
		Description: "Filesystem ACLs live outside the mode bits that `chmod` / `ls -l` / " +
			"`stat -c %a` surface. Granting `u:nobody:rwx` gives the daemon-fallback account " +
			"write access to a file; `o::rwx` / `o::rw` world-writes via ACL even when the mode " +
			"bits still look safe. Review scripts that check `stat -c %a` miss both. Prefer " +
			"`chmod` for world perms, and for specific users name the real account with the " +
			"minimum perm set.",
		Check: checkZC1595,
	})
}

func checkZC1595(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setfacl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "u:nobody:") ||
			strings.HasPrefix(v, "g:nobody:") ||
			strings.HasPrefix(v, "u:nogroup:") ||
			strings.HasPrefix(v, "g:nogroup:") ||
			v == "o::rwx" || v == "o::rw" || v == "o::rwX" {
			return []Violation{{
				KataID: "ZC1595",
				Message: "`setfacl -m " + v + "` grants perms via ACL, bypassing `chmod` / " +
					"`stat -c %a` checks. Prefer chmod for world perms, and for specific " +
					"users name the real account with minimum perms.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
