// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1740",
		Title:    "Warn on `gh release upload --clobber` — silent overwrite of release asset",
		Severity: SeverityWarning,
		Description: "`gh release upload TAG FILE --clobber` replaces an existing asset with the " +
			"same name without prompting. In production this is how a release artifact " +
			"gets silently downgraded — a CI job re-runs with a stale build and the user-" +
			"facing download moves backward without anyone noticing. Drop `--clobber` so " +
			"the second upload errors out, or version the asset name (`mytool-1.2.3-linux." +
			"tar.gz` instead of `mytool-linux.tar.gz`) so each upload has a unique slot.",
		Check: checkZC1740,
	})
}

func checkZC1740(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "release" || cmd.Arguments[1].String() != "upload" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--clobber" {
			return []Violation{{
				KataID: "ZC1740",
				Message: "`gh release upload --clobber` silently replaces an existing " +
					"asset — a re-run can downgrade the user-facing download. Drop " +
					"`--clobber` or version the asset name so each upload has a " +
					"unique slot.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
