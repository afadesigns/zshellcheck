package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1136",
		Title: "Avoid `rm -rf` without safeguard",
		Description: "`rm -rf` with a variable path is dangerous if the variable is empty. " +
			"Always validate the path or use `${var:?}` to fail on empty values.",
		Severity: SeverityWarning,
		Check:    checkZC1136,
	})
}

func checkZC1136(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "rm" {
		return nil
	}

	hasRecursiveForce := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-rf" || val == "-fr" {
			hasRecursiveForce = true
			break
		}
	}

	if !hasRecursiveForce {
		return nil
	}

	// Check if any argument is a bare variable (unprotected)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '$' {
			return []Violation{{
				KataID: "ZC1136",
				Message: "Avoid `rm -rf $var` without safeguards. Use `rm -rf ${var:?}` " +
					"to abort if the variable is empty, preventing accidental deletion.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
