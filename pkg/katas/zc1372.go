// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1372",
		Title:    "Use Zsh `zmv` autoload function instead of `rename`/`rename.ul`",
		Severity: SeverityStyle,
		Description: "Zsh's `zmv` (autoloaded via `autoload -Uz zmv`) batch-renames files using " +
			"glob patterns with capture groups. Safer than the various `rename`/`rename.ul`/`prename` " +
			"utilities (perl-based vs util-linux) and does not depend on which one is installed.",
		Check: checkZC1372,
	})
}

func checkZC1372(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "rename" && ident.Value != "rename.ul" && ident.Value != "prename" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1372",
		Message: "Use Zsh `zmv` (autoload -Uz zmv) instead of `rename`/`rename.ul`/`prename`. " +
			"Glob-pattern renaming is handled in-shell with capture groups.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
