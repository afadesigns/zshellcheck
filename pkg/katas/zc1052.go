package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1052",
		Title:       "Avoid `sed -i` for portability",
		Description: "`sed -i` usage varies between GNU/Linux and macOS/BSD. " +
			"macOS requires an extension argument (e.g. `sed -i ''`), while GNU does not. " +
			"Use a temporary file and `mv`, or `perl -i`, for portability.",
		Check: checkZC1052,
	})
}

func checkZC1052(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); !ok || name.Value != "sed" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		// Check for -i flag
		// Arguments can be PrefixExpression (-i), or Concatenated, or StringLiteral if quoted?
		// Parser usually handles `-i` as PrefixExpression(Operator="-", Right=Identifier("i"))
		// Or if it's `-i.bak`, it might be different.

		// Let's check string representation for simplicity, or specific types if robust.
		// `-i` string rep is `-i`.

		// Wait, `parsePrefixExpression` handles `-`.
		// `-i` -> Prefix(-, Ident(i)).

		if prefix, ok := arg.(*ast.PrefixExpression); ok && prefix.Operator == "-" {
			if ident, ok := prefix.Right.(*ast.Identifier); ok && ident.Value == "i" {
				violations = append(violations, Violation{
					KataID:  "ZC1052",
					Message: "`sed -i` is non-portable (GNU vs BSD differences). Use `perl -i` or a temporary file.",
					Line:    prefix.Token.Line,
					Column:  prefix.Token.Column,
				})
			}
		} else if str, ok := arg.(*ast.StringLiteral); ok {
			val := str.Value
			if len(val) >= 2 && (val[0] == '"' || val[0] == '\'') {
				val = val[1 : len(val)-1]
			}
			if val == "-i" {
				violations = append(violations, Violation{
					KataID:  "ZC1052",
					Message: "`sed -i` is non-portable (GNU vs BSD differences). Use `perl -i` or a temporary file.",
					Line:    str.Token.Line,
					Column:  str.Token.Column,
				})
			}
		}
	}

	return violations
}
