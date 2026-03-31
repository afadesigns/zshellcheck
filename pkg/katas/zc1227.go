package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1227",
		Title:    "Use `curl -f` to fail on HTTP errors",
		Severity: SeverityWarning,
		Description: "`curl` without `-f` silently returns error pages (404, 500) as success. " +
			"Use `-f` or `--fail` to return exit code 22 on HTTP errors.",
		Check: checkZC1227,
	})
}

func checkZC1227(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "curl" {
		return nil
	}

	hasFail := false
	hasURL := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-f" || val == "-fsSL" || val == "-fSL" || val == "-fsS" {
			hasFail = true
		}
		if len(val) > 4 && (val[:4] == "http" || val[:5] == "https") {
			hasURL = true
		}
	}

	if hasURL && !hasFail {
		return []Violation{{
			KataID: "ZC1227",
			Message: "Use `curl -f` to fail on HTTP errors. Without `-f`, curl silently " +
				"returns error pages (404, 500) as if they were successful.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
