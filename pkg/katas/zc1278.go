// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// Issue #343: ZC1278 fires on the same input as the canonical
// ZC1009 with overlapping advice. Stubbed to a no-op so legacy
// `disabled_katas` lists keep parsing; the detection now lives in
// ZC1009.

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1278",
		Title:       "Superseded by ZC1009 — retired duplicate",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/343 for context; the canonical detection lives in ZC1009.",
		Check:       checkZC1278,
	})
}

func checkZC1278(ast.Node) []Violation {
	return nil
}
