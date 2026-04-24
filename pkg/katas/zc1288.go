package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.DeclarationStatementNode, Kata{
		ID:       "ZC1288",
		Title:    "Use `typeset` instead of `declare` in Zsh scripts",
		Severity: SeverityStyle,
		Description: "`typeset` is the native Zsh builtin for variable declarations. " +
			"`declare` is a Bash compatibility alias. Using `typeset` is more idiomatic " +
			"and signals that the script is Zsh-native.",
		Check: checkZC1288,
		Fix:   fixZC1288,
	})
}

// fixZC1288 rewrites the `declare` keyword to `typeset`. Arguments,
// flags and assignments carry over unchanged because the two
// builtins share the same Zsh interface.
func fixZC1288(node ast.Node, v Violation, source []byte) []FixEdit {
	decl, ok := node.(*ast.DeclarationStatement)
	if !ok || decl.Command != "declare" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("declare"),
		Replace: "typeset",
	}}
}

func checkZC1288(node ast.Node) []Violation {
	decl, ok := node.(*ast.DeclarationStatement)
	if !ok {
		return nil
	}

	if decl.Command != "declare" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1288",
		Message: "Use `typeset` instead of `declare` in Zsh scripts. `typeset` is the native Zsh idiom.",
		Line:    decl.Token.Line,
		Column:  decl.Token.Column,
		Level:   SeverityStyle,
	}}
}
