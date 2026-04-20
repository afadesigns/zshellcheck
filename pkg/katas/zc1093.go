package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// Issue #341: ZC1093 fires on the same input as the canonical
// ZC1038 with overlapping advice. Stubbed to a no-op so legacy
// `disabled_katas` lists keep parsing; the detection now lives in
// ZC1038.

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1093",
		Title:       "Superseded by ZC1038 — retired duplicate",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/341 for context; the canonical detection lives in ZC1038.",
		Check:       checkZC1093,
	})
}

func checkZC1093(ast.Node) []Violation {
	return nil
}
