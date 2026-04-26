// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1305",
		Title:    "Avoid `$COMP_WORDS` — use `$words` in Zsh completion",
		Severity: SeverityWarning,
		Description: "`$COMP_WORDS` is a Bash completion variable containing the words on " +
			"the command line. Zsh completion uses `$words` array for the same purpose.",
		Check: checkZC1305,
		Fix:   fixZC1305,
	})
}

// fixZC1305 renames the Bash `$COMP_WORDS` identifier to the Zsh
// `$words` completion array. Handles both dollar-prefixed and bare
// forms.
func fixZC1305(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$COMP_WORDS":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$COMP_WORDS"),
			Replace: "$words",
		}}
	case "COMP_WORDS":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("COMP_WORDS"),
			Replace: "words",
		}}
	}
	return nil
}

func checkZC1305(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$COMP_WORDS" && ident.Value != "COMP_WORDS" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1305",
		Message: "Avoid `$COMP_WORDS` in Zsh — use `$words` array instead. `COMP_WORDS` is Bash completion-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
