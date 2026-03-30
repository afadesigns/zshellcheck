package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1153",
		Title:    "Use `cmp -s` instead of `diff` for equality check",
		Severity: SeverityStyle,
		Description: "When only checking if two files are identical (not viewing differences), " +
			"`cmp -s` is faster than `diff` as it stops at the first difference.",
		Check: checkZC1153,
	})
}

func checkZC1153(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "diff" {
		return nil
	}

	// Only flag diff -q (quiet) which is used for equality checks
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-q" {
			return []Violation{{
				KataID: "ZC1153",
				Message: "Use `cmp -s file1 file2` instead of `diff -q`. " +
					"`cmp -s` is faster for equality checks as it stops at the first difference.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
