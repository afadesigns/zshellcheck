// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1544",
		Title:    "Warn on `dnf copr enable` / `add-apt-repository ppa:` — unvetted third-party repo",
		Severity: SeverityWarning,
		Description: "Enabling a COPR project or an Ubuntu PPA pulls packages signed by a single " +
			"community contributor — there is no distro security team or reproducible-builds " +
			"guarantee behind that key. Any future compromise of that contributor's account " +
			"ships a rootkit to every box that ran this line. If you need the package badly " +
			"enough, pin to a specific `build-id`, verify the key fingerprint out of band, " +
			"and mirror to an internal repository.",
		Check: checkZC1544,
	})
}

func checkZC1544(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "dnf" && len(cmd.Arguments) >= 2 &&
		cmd.Arguments[0].String() == "copr" && cmd.Arguments[1].String() == "enable" {
		return zc1544Violation(cmd, "dnf copr enable")
	}
	if ident.Value == "add-apt-repository" {
		return zc1544Violation(cmd, "add-apt-repository")
	}
	return nil
}

func zc1544Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1544",
		Message: "`" + what + "` pulls from a single-contributor repo — no distro security " +
			"team. Pin the build, verify key fingerprint, mirror internally.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
