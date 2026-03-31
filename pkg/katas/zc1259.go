package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1259",
		Title:    "Use `grep -I` to skip binary files",
		Severity: SeverityStyle,
		Description: "`grep` without `-I` may match binary files, producing garbled output. " +
			"Use `-I` to automatically skip binary files during search.",
		Check: checkZC1259,
	})
}

func checkZC1259(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	hasSkipBinary := false
	hasRecursive := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-I" {
			hasSkipBinary = true
		}
		if val == "-r" || val == "-R" || val == "-rn" || val == "-Rn" {
			hasRecursive = true
		}
	}

	if hasRecursive && !hasSkipBinary {
		return []Violation{{
			KataID: "ZC1259",
			Message: "Use `grep -I` with recursive search to skip binary files. " +
				"Without `-I`, grep may produce garbled binary output.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
