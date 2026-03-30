package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1138",
		Title: "Use `print -l` for printing array elements one per line",
		Description: "Zsh `print -l` prints each argument on a separate line, replacing " +
			"`printf '%s\\n' \"${array[@]}\"` with a simpler, faster builtin call.",
		Severity: SeverityStyle,
		Check:    checkZC1138,
	})
}

func checkZC1138(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "printf" {
		return nil
	}

	if len(cmd.Arguments) < 1 {
		return nil
	}

	fmtStr := cmd.Arguments[0].String()
	// Match printf '%s\n' or printf "%s\n"
	clean := fmtStr
	if len(clean) >= 2 && (clean[0] == '\'' || clean[0] == '"') {
		clean = clean[1 : len(clean)-1]
	}

	if clean == "%s\\n" || clean == `%s\n` {
		return []Violation{{
			KataID: "ZC1138",
			Message: "Use `print -l` instead of `printf '%s\\n'` for printing elements one per line. " +
				"`print -l` is a Zsh builtin optimized for this pattern.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
