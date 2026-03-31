package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1255",
		Title:    "Use `curl -L` to follow HTTP redirects",
		Severity: SeverityInfo,
		Description: "`curl` without `-L` does not follow redirects, returning 301/302 responses " +
			"instead of the actual content. Use `-L` to follow redirects automatically.",
		Check: checkZC1255,
	})
}

func checkZC1255(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "curl" {
		return nil
	}

	hasFollow := false
	hasURL := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-L" || val == "-fsSL" || val == "-fSL" || val == "-sL" {
			hasFollow = true
		}
		if len(val) > 7 && val[:5] == "https" {
			hasURL = true
		}
	}

	if hasURL && !hasFollow {
		return []Violation{{
			KataID: "ZC1255",
			Message: "Use `curl -L` to follow HTTP redirects. Without `-L`, curl returns " +
				"redirect responses (301/302) instead of the actual content.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}
