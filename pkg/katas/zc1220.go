package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1220",
		Title:    "Use `chown :group` instead of `chgrp` for group changes",
		Severity: SeverityStyle,
		Description: "`chgrp` is redundant when `chown :group file` does the same thing. " +
			"Using `chown` for both user and group changes is more consistent.",
		Check: checkZC1220,
	})
}

func checkZC1220(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chgrp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1220",
		Message: "Use `chown :group file` instead of `chgrp group file`. " +
			"`chown` handles both user and group changes consistently.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
