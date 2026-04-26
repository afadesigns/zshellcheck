// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1541",
		Title:    "Error on `apk add --allow-untrusted` — installs unsigned Alpine packages",
		Severity: SeverityError,
		Description: "`apk add --allow-untrusted` skips signature verification on the package " +
			"being installed. On Alpine that is a direct MITM-to-root path: any mirror, " +
			"cache, or typo-squat can slip a replacement `.apk` and the daemon starts running " +
			"attacker code on next restart. Sign internal packages with your own key in " +
			"`/etc/apk/keys/` and keep verification on.",
		Check: checkZC1541,
	})
}

func checkZC1541(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "apk" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--allow-untrusted" {
			return []Violation{{
				KataID: "ZC1541",
				Message: "`apk --allow-untrusted` skips signature verification on the " +
					"package — MITM-to-root on Alpine. Sign and place key in /etc/apk/keys/.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
