package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1231",
		Title:    "Use `git clone --depth 1` for CI and build scripts",
		Severity: SeverityStyle,
		Description: "`git clone` without `--depth` downloads the entire history. " +
			"Use `--depth 1` in CI/build scripts where only the latest commit is needed.",
		Check: checkZC1231,
	})
}

func checkZC1231(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	if cmd.Arguments[0].String() != "clone" {
		return nil
	}

	hasDepth := false
	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--depth" || val == "--shallow-since" || val == "--single-branch" {
			hasDepth = true
		}
	}

	if !hasDepth {
		return []Violation{{
			KataID: "ZC1231",
			Message: "Consider `git clone --depth 1` in scripts. Full clones download " +
				"entire history which is unnecessary for builds and CI.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
