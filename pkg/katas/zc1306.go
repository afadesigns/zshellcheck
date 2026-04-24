package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1306",
		Title:    "Avoid `$COMP_CWORD` — use `$CURRENT` in Zsh completion",
		Severity: SeverityWarning,
		Description: "`$COMP_CWORD` is a Bash completion variable for the current cursor " +
			"word index. Zsh completion uses `$CURRENT` for the same purpose.",
		Check: checkZC1306,
		Fix:   fixZC1306,
	})
}

// fixZC1306 renames the Bash `$COMP_CWORD` identifier to the Zsh
// `$CURRENT` completion variable.
func fixZC1306(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "$COMP_CWORD":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$COMP_CWORD"),
			Replace: "$CURRENT",
		}}
	case "COMP_CWORD":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("COMP_CWORD"),
			Replace: "CURRENT",
		}}
	}
	return nil
}

func checkZC1306(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$COMP_CWORD" && ident.Value != "COMP_CWORD" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1306",
		Message: "Avoid `$COMP_CWORD` in Zsh — use `$CURRENT` instead. `COMP_CWORD` is Bash completion-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
