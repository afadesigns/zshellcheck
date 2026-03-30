package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.BooleanNode, Kata{
		ID:    "ZC1130",
		Title: "Use `:` instead of `true`",
		Description: "The `:` builtin is equivalent to `true` but guaranteed to be a shell builtin. " +
			"On some systems, `true` is an external command at /usr/bin/true.",
		Severity: SeverityStyle,
		Check:    checkZC1130,
	})
}

func checkZC1130(node ast.Node) []Violation {
	b, ok := node.(*ast.Boolean)
	if !ok || !b.Value {
		return nil
	}

	return []Violation{{
		KataID: "ZC1130",
		Message: "Use `:` instead of `true`. " +
			"`:` is always a shell builtin, while `true` may be an external command.",
		Line:   b.Token.Line,
		Column: b.Token.Column,
		Level:  SeverityStyle,
	}}
}
