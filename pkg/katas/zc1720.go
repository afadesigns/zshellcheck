package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1720",
		Title:    "Use Zsh `$COLUMNS` / `$LINES` instead of `tput cols` / `tput lines`",
		Severity: SeverityStyle,
		Description: "Zsh tracks the terminal width and height in `$COLUMNS` and `$LINES`, " +
			"updated automatically on `SIGWINCH`. Reading them is a constant-time " +
			"parameter expansion, while `tput cols` / `tput lines` forks the terminfo " +
			"helper on every call. Use the parameters; reach for `tput` only for terminfo " +
			"queries Zsh does not surface as parameters.",
		Check: checkZC1720,
	})
}

func checkZC1720(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tput" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "cols" || v == "lines" {
			repl := "$COLUMNS"
			if v == "lines" {
				repl = "$LINES"
			}
			return []Violation{{
				KataID: "ZC1720",
				Message: "Use `" + repl + "` instead of `tput " + v + "` — Zsh keeps the " +
					"terminal size in parameters, no fork needed.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}
	return nil
}
