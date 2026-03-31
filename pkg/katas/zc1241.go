package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1241",
		Title:    "Use `xargs -0` with null separators for safe argument passing",
		Severity: SeverityWarning,
		Description: "`xargs` without `-0` splits on whitespace, breaking on filenames with spaces. " +
			"Use `xargs -0` paired with `find -print0` for safe handling.",
		Check: checkZC1241,
	})
}

func checkZC1241(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}

	hasNull := false
	hasRM := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-0" {
			hasNull = true
		}
		if val == "rm" {
			hasRM = true
		}
	}

	if hasRM && !hasNull {
		return []Violation{{
			KataID: "ZC1241",
			Message: "Use `xargs -0 rm` with `find -print0` for safe deletion. " +
				"Without `-0`, filenames with spaces or special characters break.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
