package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1715",
		Title:    "Error on `read -p \"prompt\"` — Zsh `-p` reads from coprocess, not a prompt",
		Severity: SeverityError,
		Description: "Bash's `read -p \"Prompt: \" var` prints the prompt before reading. " +
			"Zsh's `read -p` means \"read from the coprocess set up with `coproc`\" — when " +
			"no coprocess exists, `read` errors with `no coprocess` and leaves the variable " +
			"empty, silently breaking the script. The Zsh idiom is `read \"var?Prompt: \"` " +
			"— a `?` after the variable name introduces the prompt string, with the same " +
			"behavior under `-r`, `-s`, etc.",
		Check: checkZC1715,
	})
}

func checkZC1715(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "read" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if len(v) < 2 || v[0] != '-' {
			continue
		}
		// Skip long flags (none of read's short flags need to be considered as long).
		if v[1] == '-' {
			continue
		}
		if !strings.ContainsRune(v[1:], 'p') {
			continue
		}
		return []Violation{{
			KataID: "ZC1715",
			Message: "`read " + v + "` triggers Zsh's coprocess reader, not Bash's prompt — " +
				"the variable stays empty. Use `read \"var?Prompt: \"` (the `?` after the " +
				"variable name introduces the prompt).",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}
