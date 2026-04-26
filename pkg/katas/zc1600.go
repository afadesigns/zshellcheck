// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1600",
		Title:    "Warn on bare `chroot DIR CMD` — missing `--userspec=` keeps uid 0 inside the jail",
		Severity: SeverityWarning,
		Description: "`chroot` changes the filesystem root but does not drop privileges. The " +
			"caller is almost always root (the syscall needs `CAP_SYS_CHROOT`), and without " +
			"`--userspec=USER:GROUP` the command inside the chroot still runs as uid 0. It can " +
			"write anywhere inside the tree, chmod binaries, and — if proc / sys / device nodes " +
			"are bind-mounted in — escape. Pass `--userspec=` to run the command as a named " +
			"unprivileged user, or drop to a dedicated helper (bubblewrap, firejail) that also " +
			"unshares user namespaces.",
		Check: checkZC1600,
	})
}

func checkZC1600(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chroot" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1600",
		Message: "`chroot` without `--userspec=` runs the inner command as uid 0. Pass " +
			"`--userspec=USER:GROUP` to drop privileges, or use `bwrap` / `firejail` for " +
			"user-namespace isolation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
