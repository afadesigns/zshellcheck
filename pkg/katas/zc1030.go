package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(reflect.TypeOf(&ast.SimpleCommand{}), Kata{
		ID:          "ZC1030",
		Title:       "Use `printf` instead of `echo`",
		Description: "The `echo` command's behavior can be inconsistent across different shells and " +
			"environments, especially with flags and escape sequences. `printf` provides more reliable " +
			"and portable string formatting.",
		Check:       checkZC1030,
	})
}

func checkZC1030(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name.TokenLiteral() != "echo" {
		return nil
	}

	// Defer to ZC1037 if any argument is a variable.
	for _, arg := range cmd.Arguments {
		if ident, ok := arg.(*ast.Identifier); ok {
			if ident.Token.Type == "VARIABLE" {
				return nil
			}
		}
	}

	return []Violation{
		{
			KataID:  "ZC1030",
			Message: "Use `printf` for more reliable and portable string formatting instead of `echo`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
		},
	}
}

