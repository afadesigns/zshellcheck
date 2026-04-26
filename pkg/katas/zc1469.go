// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1469",
		Title:    "Error on `dnf/yum --nogpgcheck` or `rpm --nosignature` (unsigned RPM install)",
		Severity: SeverityError,
		Description: "`--nogpgcheck` / `--nosignature` / `--nodigest` disable RPM package " +
			"signature and digest verification. This turns every mirror, cache, or MITM into a " +
			"direct root compromise. Always keep GPG/signature checking on; sign internal repos " +
			"with your own key.",
		Check: checkZC1469,
	})
}

func checkZC1469(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "dnf", "yum", "microdnf", "zypper":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "--nogpgcheck" || v == "--no-gpg-checks" {
				return zc1469Violation(cmd, ident.Value+" "+v)
			}
		}
	case "rpm", "rpmbuild":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "--nosignature" || v == "--nodigest" || v == "--nofiledigest" ||
				v == "--noverify" || v == "--nochecksum" {
				return zc1469Violation(cmd, ident.Value+" "+v)
			}
		}
	}
	return nil
}

func zc1469Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1469",
		Message: "Package signature verification disabled (" + what + ") — any mirror / MITM " +
			"becomes immediate root.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
