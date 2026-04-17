package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1413",
		Title:    "Use Zsh `whence -p cmd` instead of `hash -t cmd` for resolved path",
		Severity: SeverityStyle,
		Description: "Bash's `hash -t cmd` prints the hashed path for `cmd` (or fails if not " +
			"hashed). Zsh's `whence -p cmd` prints the PATH-resolved absolute path, whether " +
			"hashed or not — more reliable and the native Zsh idiom.",
		Check: checkZC1413,
	})
}

func checkZC1413(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "hash" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-t" {
			return []Violation{{
				KataID: "ZC1413",
				Message: "Use `whence -p cmd` (Zsh) instead of `hash -t cmd`. " +
					"`whence -p` always returns the absolute path, regardless of hash state.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
