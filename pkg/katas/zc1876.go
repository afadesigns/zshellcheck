// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1876",
		Title:    "Warn on `cargo publish --allow-dirty` — publishes the crate with uncommitted local changes",
		Severity: SeverityWarning,
		Description: "`cargo publish` by default refuses to upload when the working tree is dirty, " +
			"because the published tarball is a snapshot of whatever is on disk — not " +
			"whatever is committed. `--allow-dirty` skips that check, so a `println!` " +
			"dropped in for debugging, an uncommitted `Cargo.toml` dep bump, or a " +
			"`patch.crates-io` override that only exists locally ends up on crates.io " +
			"under the same version users see on GitHub. This is irreversible — once a " +
			"version is uploaded it cannot be replaced, only yanked. Commit first and " +
			"publish from a clean checkout; if you truly must publish from a dirty tree, " +
			"scope the flag to a one-off manual call with a `--dry-run` rehearsal first.",
		Check: checkZC1876,
	})
}

func checkZC1876(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cargo" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}
	if args[0].String() != "publish" {
		return nil
	}
	for _, arg := range args[1:] {
		if arg.String() == "--allow-dirty" {
			return []Violation{{
				KataID: "ZC1876",
				Message: "`cargo publish --allow-dirty` uploads a tarball snapshot of " +
					"the dirty working tree — debug prints and local-only patches " +
					"end up on crates.io for a version that cannot be replaced. " +
					"Commit first; `--dry-run` to rehearse.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
