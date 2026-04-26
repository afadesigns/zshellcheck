// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// Issue #344: ZC1277 fires on the same input as the canonical
// ZC1108 with overlapping advice. Stubbed to a no-op so legacy
// `disabled_katas` lists keep parsing; the detection now lives in
// ZC1108.

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1277",
		Title:       "Superseded by ZC1108 — retired duplicate",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/344 for context; the canonical detection lives in ZC1108.",
		Check:       checkZC1277,
	})
}

func checkZC1277(ast.Node) []Violation {
	return nil
}
