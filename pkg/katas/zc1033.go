package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// Per issue #345: ZC1023–ZC1029, ZC1033, ZC1035 were identical copies of
// the canonical ZC1022 `let` detection, each emitting the same violation
// on the same input. Ten fires per `let` line inflated user-visible noise.
// Project rule ("once committed, fix — don't remove") keeps the IDs alive
// as no-op stubs so existing `disabled_katas` lists that reference them
// keep parsing, but the duplicate detection is retired in favour of
// ZC1022.

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:          "ZC1033",
		Title:       "Superseded by ZC1022 — retired duplicate `let` detector",
		Severity:    SeverityStyle,
		Description: "Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.",
		Check:       checkZC1033,
	})
}

func checkZC1033(ast.Node) []Violation {
	return nil
}
