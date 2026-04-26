// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1515",
		Title:    "Warn on `md5sum` / `sha1sum` for integrity check — collision-vulnerable",
		Severity: SeverityWarning,
		Description: "MD5 and SHA-1 are broken for collision resistance: public attacks cheaply " +
			"craft two different files with the same hash. For verifying a download against a " +
			"published checksum, or for comparing archives against a manifest, use " +
			"`sha256sum` / `sha512sum` / `b2sum` instead. MD5 is still fine for non-adversarial " +
			"cache keys but almost every invocation in scripts is the integrity case.",
		Check: checkZC1515,
	})
}

func checkZC1515(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "md5sum" && ident.Value != "sha1sum" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1515",
		Message: "`" + ident.Value + "` is collision-vulnerable — don't use for integrity " +
			"checks. Use `sha256sum` / `sha512sum` / `b2sum` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
