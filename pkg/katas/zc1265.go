package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1265",
		Title:    "Use `systemctl enable --now` to enable and start together",
		Severity: SeverityStyle,
		Description: "`systemctl enable` without `--now` only enables on next boot. " +
			"Use `--now` to enable and immediately start the service.",
		Check: checkZC1265,
	})
}

func checkZC1265(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}

	hasEnable := false
	hasNow := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "enable" {
			hasEnable = true
		}
		if val == "--now" {
			hasNow = true
		}
	}

	if hasEnable && !hasNow {
		return []Violation{{
			KataID: "ZC1265",
			Message: "Use `systemctl enable --now` to enable and start the service immediately. " +
				"Without `--now`, the service only starts on next boot.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
