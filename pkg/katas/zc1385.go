package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1385",
		Title:    "Avoid `$PS0` — Bash-only; Zsh uses `preexec` hook",
		Severity: SeverityWarning,
		Description: "Bash 4.4+ prints `$PS0` after reading a command and before executing it. Zsh " +
			"does not honor `$PS0`; the equivalent is a `preexec` function (or " +
			"`add-zsh-hook preexec funcname`) which receives the command line as `$1`.",
		Check: checkZC1385,
	})
}

func checkZC1385(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" && ident.Value != "export" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "$PS0") || strings.Contains(v, "${PS0}") || v == "PS0" {
			return []Violation{{
				KataID: "ZC1385",
				Message: "`$PS0` is Bash-only. Zsh uses the `preexec` hook function for " +
					"pre-execution prompts.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
