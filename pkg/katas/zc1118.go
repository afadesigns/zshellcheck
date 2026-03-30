package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1118",
		Title: "Use `print -rn` instead of `echo -n`",
		Description: "The behavior of `echo -n` varies across shells and platforms. " +
			"In Zsh, `print -rn` is the reliable way to output text without a trailing newline.",
		Severity: SeverityStyle,
		Check:    checkZC1118,
	})
}

func checkZC1118(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "echo" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	firstArg := cmd.Arguments[0].String()
	if firstArg == "-n" {
		return []Violation{{
			KataID: "ZC1118",
			Message: "Use `print -rn` instead of `echo -n`. " +
				"`echo -n` behavior varies across shells; `print -rn` is the reliable Zsh idiom.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
