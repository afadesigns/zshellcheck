package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1005",
		Title: "Use whence instead of which",
		Description: "The `which` command is an external command and may not be available on all systems. " +
			"The `whence` command is a built-in Zsh command that provides a more reliable and consistent " +
			"way to find the location of a command.",
		Severity: SeverityInfo,
		Check:    checkZC1005,
		Fix:      fixZC1005,
	})
}

// fixZC1005 rewrites `which` -> `whence` at the command name position.
// Arguments are unchanged; the two builtins share the identifier-query
// shape for the common case.
func fixZC1005(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "which" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("which"),
		Replace: "whence",
	}}
}

func checkZC1005(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "which" {
				violations = append(violations, Violation{
					KataID: "ZC1005",
					Message: "Use whence instead of which. The `whence` command is a built-in Zsh command " +
						"that provides a more reliable and consistent way to find the location of a command.",
					Line:   ident.Token.Line,
					Column: ident.Token.Column,
					Level:  SeverityInfo,
				})
			}
		}
	}

	return violations
}
