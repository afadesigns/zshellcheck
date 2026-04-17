package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1374",
		Title:    "Avoid `$FUNCNEST` — Zsh uses `$FUNCNEST` as a limit, not a depth indicator",
		Severity: SeverityWarning,
		Description: "Bash's `$FUNCNEST` is both a writable limit and (implicitly) the current " +
			"depth-query vehicle. Zsh's `$FUNCNEST` is only the limit — to read the current depth " +
			"use `${#funcstack}`. Reading `$FUNCNEST` expecting depth returns the limit, not " +
			"the current depth.",
		Check: checkZC1374,
	})
}

func checkZC1374(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "$FUNCNEST" || v == "${FUNCNEST}" {
			return []Violation{{
				KataID: "ZC1374",
				Message: "In Zsh, `$FUNCNEST` is the configured limit, not the current depth. " +
					"Use `${#funcstack}` for current function nesting depth.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
