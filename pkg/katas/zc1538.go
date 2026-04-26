// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1538",
		Title:    "Error on `zpool destroy -f` / `zfs destroy -rR` — recursive ZFS destruction",
		Severity: SeverityError,
		Description: "`zpool destroy -f` nukes a whole ZFS pool including every dataset, " +
			"snapshot, and clone on it. `zfs destroy -r` recurses into descendant datasets; " +
			"`-R` additionally drops descendant clones. Unlike `rm`, the space is freed " +
			"immediately and there is no recycle bin. Always require `zfs list`/`zpool list` " +
			"+ explicit target confirmation in the same script block, and prefer snapshot-" +
			"based rollback for recoverable workflows.",
		Check: checkZC1538,
	})
}

func checkZC1538(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "zpool" && len(cmd.Arguments) >= 2 &&
		cmd.Arguments[0].String() == "destroy" {
		for _, arg := range cmd.Arguments[1:] {
			if arg.String() == "-f" {
				return zc1538Violation(cmd, "zpool destroy -f")
			}
		}
	}
	if ident.Value == "zfs" && len(cmd.Arguments) >= 2 &&
		cmd.Arguments[0].String() == "destroy" {
		for _, arg := range cmd.Arguments[1:] {
			v := arg.String()
			if v == "-r" || v == "-R" || v == "-rR" || v == "-Rr" ||
				v == "-rf" || v == "-fr" {
				return zc1538Violation(cmd, "zfs destroy "+v)
			}
		}
	}
	return nil
}

func zc1538Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1538",
		Message: "`" + what + "` irrecoverably destroys the ZFS pool/dataset and every " +
			"snapshot on it. Require explicit target confirmation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
