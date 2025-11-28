package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1076",
		Title: "Use `autoload -Uz` for lazy loading",
		Description: "When using `autoload`, prefer `-Uz` to ensure standard Zsh behavior (no alias expansion, zsh style). " +
			"`-U` prevents alias expansion, and `-z` ensures Zsh style autoloading.",
		Check: checkZC1076,
	})
}

func checkZC1076(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name.String() != "autoload" {
		return nil
	}

	hasU := false
	hasZ := false

	for _, arg := range cmd.Arguments {
		argStr := arg.String()
		argStr = strings.Trim(argStr, "\"'")
		if strings.HasPrefix(argStr, "-") {
			if strings.Contains(argStr, "U") {
				hasU = true
			}
			if strings.Contains(argStr, "z") {
				hasZ = true
			}
		}
	}

	if !hasU || !hasZ {
		return []Violation{{
			KataID:  "ZC1076",
			Message: "Use `autoload -Uz` to ensure consistent and safe function loading.",
			Line:    cmd.TokenLiteralNode().Line,
			Column:  cmd.TokenLiteralNode().Column,
		}}
	}

	return nil
}
