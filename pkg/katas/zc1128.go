package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1128",
		Title: "Use `> file` instead of `touch file` for creation",
		Description: "If the goal is to create an empty file, `> file` does it without " +
			"spawning `touch`. Use `touch` only when you need to update timestamps.",
		Severity: SeverityStyle,
		Check:    checkZC1128,
	})
}

func checkZC1128(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "touch" {
		return nil
	}

	// Skip touch with flags (timestamps: -t, -d, -r, -a, -m)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	// Only flag touch with a single file argument
	if len(cmd.Arguments) != 1 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1128",
		Message: "Use `> file` instead of `touch file` to create an empty file. " +
			"This avoids spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
