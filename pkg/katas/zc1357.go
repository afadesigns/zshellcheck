package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1357",
		Title:    "Use Zsh `${(q)var}` instead of `printf '%q'` for shell-quoting",
		Severity: SeverityStyle,
		Description: "Bash's `printf '%q'` emits shell-quoted output. Zsh's `${(q)var}` parameter " +
			"flag does the same in-shell, with variants `${(qq)var}`, `${(qqq)var}`, `${(qqqq)var}` " +
			"for single-quote, double-quote, $'...', and POSIX ANSI-C styles respectively.",
		Check: checkZC1357,
	})
}

func checkZC1357(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		// Any format string containing %q (unescaped) triggers the kata.
		if strings.Contains(val, "%q") {
			return []Violation{{
				KataID: "ZC1357",
				Message: "Use Zsh `${(q)var}` for shell-quoting instead of `printf '%q'`. " +
					"Variants: `${(qq)}`, `${(qqq)}`, `${(qqqq)}` for different quote styles.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
