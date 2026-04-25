package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1297",
		Title:    "Avoid `$BASH_SOURCE` — use `$0` or `${(%):-%x}` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_SOURCE` is a Bash-specific variable that does not exist in Zsh. " +
			"In Zsh, use `$0` inside a sourced file to get the script path, or " +
			"`${(%):-%x}` for the current file regardless of sourcing context.",
		Check: checkZC1297,
		Fix:   fixZC1297,
	})
}

// fixZC1297 renames the Bash `$BASH_SOURCE` identifier to the Zsh
// `${(%):-%x}` prompt-flag expansion that resolves to the current file
// regardless of sourcing context. Only the dollar-prefixed form is
// rewritten — the bare `BASH_SOURCE` form (inside `${...}` or as an
// assignment target) is left to manual review because the surrounding
// braces would need adjusting too.
func fixZC1297(node ast.Node, v Violation, _ []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok || ident == nil {
		return nil
	}
	if ident.Value != "$BASH_SOURCE" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("$BASH_SOURCE"),
		Replace: "${(%):-%x}",
	}}
}

func checkZC1297(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_SOURCE" && ident.Value != "BASH_SOURCE" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1297",
		Message: "Avoid `$BASH_SOURCE` in Zsh — use `$0` or `${(%):-%x}` instead. `BASH_SOURCE` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
