// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1968",
		Title:    "Warn on `dnf versionlock add` / `yum versionlock add` — pins RPM, blocks CVE updates",
		Severity: SeverityWarning,
		Description: "`dnf versionlock add pkg` (and the legacy `yum versionlock add pkg`) " +
			"write an entry to `/etc/dnf/plugins/versionlock.list` that excludes the " +
			"package from future `dnf update` / `dnf upgrade` runs. Mirrors `apt-mark " +
			"hold` on Debian (ZC1550): the lock survives reboots and unattended-upgrades " +
			"never sees the newer rpm, so kernel, openssl, or glibc CVEs pile up unseen. " +
			"Document the exact reason in the commit, pair the lock with a scheduled " +
			"`dnf versionlock delete` date, and prefer excluding the problematic " +
			"transaction via `--exclude` or a one-shot `dnf update --setopt=exclude=pkg` " +
			"rather than a persistent pin.",
		Check: checkZC1968,
	})
}

func checkZC1968(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dnf" && ident.Value != "yum" && ident.Value != "microdnf" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "versionlock" {
		return nil
	}
	sub := cmd.Arguments[1].String()
	if sub != "add" && sub != "exclude" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1968",
		Message: "`" + ident.Value + " versionlock " + sub + "` pins the rpm — blocks " +
			"future CVE fixes for glibc/openssl/kernel. Prefer `--exclude` on a single " +
			"transaction and schedule a `versionlock delete` review.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
