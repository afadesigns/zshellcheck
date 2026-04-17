package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1355",
		Title:    "Use `print -r` instead of `echo -E` for raw output",
		Severity: SeverityStyle,
		Description: "`echo -E` disables backslash interpretation, but the flag is Bash-ism and " +
			"ignored by POSIX `echo`. Zsh's `print -r` is the idiomatic raw-printer; combine " +
			"with `-n` (no newline), `-l` (one per line), `-u<fd>` (file descriptor), or `--` " +
			"(end of flags) as needed.",
		Check: checkZC1355,
	})
}

func checkZC1355(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "echo" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-E" {
			return []Violation{{
				KataID: "ZC1355",
				Message: "Use `print -r` instead of `echo -E` for raw output. " +
					"`-E` is a Bash-ism and ignored by POSIX echo.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
