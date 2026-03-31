package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1268",
		Title:    "Use `du -sh --` to handle filenames starting with dash",
		Severity: SeverityInfo,
		Description: "`du -sh *` breaks if a filename starts with `-`. " +
			"Use `--` to signal end of options and safely handle all filenames.",
		Check: checkZC1268,
	})
}

func checkZC1268(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "du" {
		return nil
	}

	hasEndOfOpts := false
	hasGlob := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "--" {
			hasEndOfOpts = true
		}
		if val == "*" || val == "." {
			hasGlob = true
		}
	}

	if hasGlob && !hasEndOfOpts {
		return []Violation{{
			KataID: "ZC1268",
			Message: "Use `du -sh -- *` instead of `du -sh *`. The `--` prevents " +
				"filenames starting with `-` from being interpreted as options.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}
