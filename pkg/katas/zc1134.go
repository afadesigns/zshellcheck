package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1134",
		Title: "Avoid `sleep` in tight loops",
		Description: "Using `sleep` inside a loop for polling creates busy-wait patterns. " +
			"Consider `inotifywait`, `zle`, or event-driven approaches instead.",
		Severity: SeverityStyle,
		Check:    checkZC1134,
	})
}

func checkZC1134(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sleep" {
		return nil
	}

	// Flag sleep with very short intervals (0, 0.1, 0.5, 1)
	if len(cmd.Arguments) != 1 {
		return nil
	}

	val := cmd.Arguments[0].String()
	if val == "0" || val == "0.1" || val == "0.01" || val == "0.5" {
		return []Violation{{
			KataID: "ZC1134",
			Message: "Avoid `sleep " + val + "` in loops. Short sleep intervals suggest busy-waiting. " +
				"Consider event-driven alternatives like `inotifywait` or `zle -F`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
