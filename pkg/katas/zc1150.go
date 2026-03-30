package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1150",
		Title:    "Avoid `cat` with single file argument in assignment",
		Severity: SeverityStyle,
		Description: "Use `$(<file)` instead of `$(cat file)` to read file contents. " +
			"Zsh's `$(<file)` is a built-in that avoids spawning cat.",
		Check: checkZC1150,
	})
}

func checkZC1150(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	// Flag cat with exactly one non-flag argument (likely $(cat file) pattern)
	fileArgs := 0
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] != '-' {
			fileArgs++
		} else {
			return nil // has flags, not simple cat
		}
	}

	if fileArgs == 1 {
		return []Violation{{
			KataID: "ZC1150",
			Message: "Use `$(<file)` instead of `$(cat file)` to read file contents. " +
				"Zsh's `$(<file)` is a built-in that avoids spawning cat.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
