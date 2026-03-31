package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1253",
		Title:    "Use `docker build --no-cache` in CI for reproducible builds",
		Severity: SeverityStyle,
		Description: "`docker build` uses layer caching which can mask dependency changes. " +
			"Use `--no-cache` in CI pipelines to ensure fully reproducible builds.",
		Check: checkZC1253,
	})
}

func checkZC1253(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "build" {
		return nil
	}

	hasNoCache := false
	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--no-cache" {
			hasNoCache = true
		}
	}

	if !hasNoCache {
		return []Violation{{
			KataID: "ZC1253",
			Message: "Consider `docker build --no-cache` in CI for reproducible builds. " +
				"Layer caching can mask changed dependencies.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
