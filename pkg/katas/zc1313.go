package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1313",
		Title:    "Avoid `$BASH_ALIASES` — use Zsh `aliases` hash",
		Severity: SeverityWarning,
		Description: "`$BASH_ALIASES` is a Bash associative array of defined aliases. " +
			"Zsh provides the `aliases` associative array for the same purpose.",
		Check: checkZC1313,
		Fix:   fixZC1313,
	})
}

// fixZC1313 renames the Bash `$BASH_ALIASES` identifier to the Zsh
// `$aliases` associative array.
func fixZC1313(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "$BASH_ALIASES":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$BASH_ALIASES"),
			Replace: "$aliases",
		}}
	case "BASH_ALIASES":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("BASH_ALIASES"),
			Replace: "aliases",
		}}
	}
	return nil
}

func checkZC1313(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$BASH_ALIASES" && ident.Value != "BASH_ALIASES" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1313",
		Message: "Avoid `$BASH_ALIASES` in Zsh — use the `aliases` associative array instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
