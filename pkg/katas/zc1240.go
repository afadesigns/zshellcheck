package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1240",
		Title:    "Use `find -maxdepth` with `-delete` to limit scope",
		Severity: SeverityWarning,
		Description: "`find -delete` without `-maxdepth` recurses infinitely and may " +
			"delete more than intended. Always limit the search depth.",
		Check: checkZC1240,
	})
}

func checkZC1240(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	hasDelete := false
	hasMaxdepth := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-delete" {
			hasDelete = true
		}
		if val == "-maxdepth" {
			hasMaxdepth = true
		}
	}

	if hasDelete && !hasMaxdepth {
		return []Violation{{
			KataID: "ZC1240",
			Message: "Use `find -maxdepth N` with `-delete` to limit deletion scope. " +
				"Without depth limits, find recurses infinitely.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
