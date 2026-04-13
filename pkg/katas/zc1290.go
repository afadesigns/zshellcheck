package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1290",
		Title:    "Use Zsh `${(n)array}` for numeric sorting instead of `sort -n`",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(n)` parameter expansion flag to sort array elements " +
			"numerically. This avoids spawning an external `sort -n` process for " +
			"simple numeric sorting of array data.",
		Check: checkZC1290,
	})
}

func checkZC1290(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sort" {
		return nil
	}

	hasNumeric := false
	hasOtherFlags := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-n" {
			hasNumeric = true
		} else if len(val) > 1 && val[0] == '-' {
			hasOtherFlags = true
		}
	}

	if hasNumeric && !hasOtherFlags {
		return []Violation{{
			KataID:  "ZC1290",
			Message: "Use Zsh `${(n)array}` for numeric sorting instead of `sort -n`. The `(n)` flag sorts numerically in-shell.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}
