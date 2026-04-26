// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1298",
		Title:    "Avoid `$FUNCNAME` — use `$funcstack` in Zsh",
		Severity: SeverityWarning,
		Description: "`$FUNCNAME` is a Bash-specific array that does not exist in Zsh. " +
			"Zsh provides `$funcstack` as the equivalent, containing the call stack " +
			"of function names with the current function at index 1.",
		Check: checkZC1298,
		Fix:   fixZC1298,
	})
}

// fixZC1298 renames the Bash `$FUNCNAME` identifier to the Zsh
// `$funcstack` equivalent. Handles both the dollar-prefixed and
// bare forms.
func fixZC1298(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$FUNCNAME":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$FUNCNAME"),
			Replace: "$funcstack",
		}}
	case "FUNCNAME":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("FUNCNAME"),
			Replace: "funcstack",
		}}
	}
	return nil
}

func checkZC1298(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$FUNCNAME" && ident.Value != "FUNCNAME" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1298",
		Message: "Avoid `$FUNCNAME` in Zsh — use `$funcstack` instead. `FUNCNAME` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
