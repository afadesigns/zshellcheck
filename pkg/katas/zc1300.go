package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1300",
		Title:    "Avoid `$BASH_VERSINFO` — use `$ZSH_VERSION` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_VERSINFO` is a Bash-specific array containing version components. " +
			"In Zsh, use `$ZSH_VERSION` (string) or `${(s:.:)ZSH_VERSION}` to split " +
			"it into components for version comparison.",
		Check: checkZC1300,
		Fix:   fixZC1300,
	})
}

// fixZC1300 renames `$BASH_VERSION` / `$BASH_VERSINFO` to the Zsh
// equivalent `$ZSH_VERSION`. The lossy case (BASH_VERSINFO is an
// array, ZSH_VERSION is a string) is the best single-token swap
// available; callers that need components can split the string with
// the `${(s:.:)ZSH_VERSION}` flag.
func fixZC1300(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "$BASH_VERSION", "$BASH_VERSINFO":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len(ident.Value),
			Replace: "$ZSH_VERSION",
		}}
	case "BASH_VERSION", "BASH_VERSINFO":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len(ident.Value),
			Replace: "ZSH_VERSION",
		}}
	}
	return nil
}

func checkZC1300(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$BASH_VERSINFO" && ident.Value != "BASH_VERSINFO" &&
		ident.Value != "$BASH_VERSION" && ident.Value != "BASH_VERSION" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1300",
		Message: "Avoid Bash version variables in Zsh — use `$ZSH_VERSION` instead. Bash version variables are undefined in Zsh.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
