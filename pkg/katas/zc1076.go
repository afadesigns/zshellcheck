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
		// The parser might treat flags as PrefixExpression or Identifier depending on context
		// But usually simple flags are handled as part of arguments list if they are simple strings
		
		// Check if it's a prefix expression (e.g. -Uz)
		if pe, ok := arg.(*ast.PrefixExpression); ok && pe.Operator == "-" {
			if ident, ok := pe.Right.(*ast.Identifier); ok {
				if strings.Contains(ident.Value, "U") {
					hasU = true
				}
				if strings.Contains(ident.Value, "z") {
					hasZ = true
				}
			}
		}
		// Or if parser just treats it as Identifier starting with -
		// (This depends on lexer/parser specifics, usually flags are parsed as PrefixExpr or just Identifier if space separated?)
		// Let's assume the test failure means it wasn't catching the string check properly in previous attempt because of AST structure.
		// Previous attempt checked `arg.String()` which might reconstruct `-` + `U` but AST might be PrefixExpr.
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