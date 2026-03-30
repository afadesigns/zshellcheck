package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1170",
		Title:    "Avoid `pushd`/`popd` without `-q` flag",
		Severity: SeverityStyle,
		Description: "`pushd` and `popd` print the directory stack by default, cluttering output. " +
			"Use `-q` flag to suppress output in scripts.",
		Check: checkZC1170,
	})
}

func checkZC1170(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "pushd" && ident.Value != "popd" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-q" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1170",
		Message: "Use `" + ident.Value + " -q` to suppress directory stack output in scripts. " +
			"Without `-q`, the stack is printed on every call.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
