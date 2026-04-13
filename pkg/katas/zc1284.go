package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1284",
		Title:    "Use Zsh `${(s:sep:)var}` instead of `cut -d` for field splitting",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(s:separator:)` parameter expansion flag to split strings " +
			"into arrays by a delimiter. This is more idiomatic than invoking `cut -d` and " +
			"avoids spawning an external process.",
		Check: checkZC1284,
	})
}

func checkZC1284(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cut" {
		return nil
	}

	hasDelim := false
	hasField := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if strings.HasPrefix(val, "-d") && val != "-d." {
			hasDelim = true
		}
		if strings.HasPrefix(val, "-f") {
			hasField = true
		}
	}

	if hasDelim && hasField {
		return []Violation{{
			KataID:  "ZC1284",
			Message: "Use Zsh parameter expansion `${(s:sep:)var}` for field splitting instead of `cut -d -f`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}
