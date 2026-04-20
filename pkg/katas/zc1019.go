package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// Issue #342: ZC1019 fires on the same input as the canonical
// ZC1005 with overlapping advice. Stubbed to a no-op so legacy
// `disabled_katas` lists keep parsing; the detection now lives in
// ZC1005.

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1019",
		Title:       "Superseded by ZC1005 — retired duplicate",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/342 for context; the canonical detection lives in ZC1005.",
		Check:       checkZC1019,
	})
}

func checkZC1019(ast.Node) []Violation {
	return nil
}
