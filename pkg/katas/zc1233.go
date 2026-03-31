package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1233",
		Title:    "Avoid `npm install -g` — use `npx` for one-off tools",
		Severity: SeverityStyle,
		Description: "Global npm installs pollute the system. Use `npx` to run tools " +
			"without installing, or `npm install --save-dev` for project dependencies.",
		Check: checkZC1233,
	})
}

func checkZC1233(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "npm" {
		return nil
	}

	hasInstall := false
	hasGlobal := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "install" || val == "i" {
			hasInstall = true
		}
		if val == "-g" {
			hasGlobal = true
		}
	}

	if hasInstall && hasGlobal {
		return []Violation{{
			KataID: "ZC1233",
			Message: "Avoid `npm install -g`. Use `npx` for one-off tool execution " +
				"or `npm install --save-dev` for project-scoped dependencies.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
