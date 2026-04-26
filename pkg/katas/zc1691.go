// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1691",
		Title:    "Warn on `rsync --remove-source-files` — SRC deletion tied to optimistic success",
		Severity: SeverityWarning,
		Description: "`rsync --remove-source-files` deletes each source file once rsync has " +
			"transferred it. The delete is gated on rsync's per-file success, which is " +
			"generous: a remote out-of-disk error after the partial write, a `--chmod` " +
			"rejection, or a flaky network that drops after the data bytes but before " +
			"metadata can still look like success. Couple that with a wrong DST path and " +
			"the source is gone with nothing to recover. Prefer a two-step flow: `rsync " +
			"-a SRC DST` first, verify DST (checksums / file count), then `rm` the source " +
			"explicitly, or use `mv` for local moves.",
		Check: checkZC1691,
	})
}

func checkZC1691(node ast.Node) []Violation {
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
		if arg.String() == "--remove-source-files" {
			return []Violation{{
				KataID: "ZC1691",
				Message: "`rsync --remove-source-files` deletes SRC on optimistic per-file " +
					"success — verify DST after the transfer and `rm` explicitly instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
