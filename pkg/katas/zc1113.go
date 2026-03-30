package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1113",
		Title: "Use `${var:A}` instead of `realpath` or `readlink -f`",
		Description: "Zsh provides the `:A` modifier to resolve a path to its absolute form, " +
			"following symlinks. Avoid spawning `realpath` or `readlink -f` as external processes.",
		Check: checkZC1113,
	})
}

func checkZC1113(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value

	if name == "realpath" {
		for _, arg := range cmd.Arguments {
			val := arg.String()
			if len(val) > 1 && val[0] == '-' && val != "-s" {
				return nil
			}
		}
		return []Violation{{
			KataID: "ZC1113",
			Message: "Use `${var:A}` instead of `realpath` to resolve absolute paths. " +
				"Zsh path modifiers avoid spawning an external process.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
		}}
	}

	if name == "readlink" {
		hasResolveFlag := false
		for _, arg := range cmd.Arguments {
			val := arg.String()
			if val == "-f" || val == "-e" || val == "-m" {
				hasResolveFlag = true
			}
		}
		if !hasResolveFlag {
			return nil
		}
		return []Violation{{
			KataID: "ZC1113",
			Message: "Use `${var:A}` instead of `readlink -f` to resolve absolute paths. " +
				"Zsh path modifiers avoid spawning an external process.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
		}}
	}

	return nil
}
