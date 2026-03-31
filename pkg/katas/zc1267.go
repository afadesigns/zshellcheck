package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1267",
		Title:    "Use `df -P` for POSIX-portable disk usage output",
		Severity: SeverityStyle,
		Description: "`df -h` output format varies across systems and locales. " +
			"Use `df -P` for single-line, fixed-format output safe for script parsing.",
		Check: checkZC1267,
	})
}

func checkZC1267(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "df" {
		return nil
	}

	hasPortable := false
	hasHuman := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-P" {
			hasPortable = true
		}
		if val == "-h" {
			hasHuman = true
		}
	}

	if hasHuman && !hasPortable {
		return []Violation{{
			KataID: "ZC1267",
			Message: "Use `df -P` for script-safe output. `df -h` format varies across " +
				"systems and may split long device names across lines.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
