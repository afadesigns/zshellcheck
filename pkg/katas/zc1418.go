package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1418",
		Title:    "Use Zsh `limit -h`/`-s` instead of `ulimit -H`/`-S` for hard/soft limits",
		Severity: SeverityStyle,
		Description: "Bash's `ulimit` uses uppercase `-H` (hard) and `-S` (soft). Zsh's native " +
			"`limit` builtin uses lowercase `-h` and `-s` for the same. The Zsh form is easier " +
			"to remember and produces human-readable output.",
		Check: checkZC1418,
	})
}

func checkZC1418(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ulimit" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-H" || v == "-S" || v == "-HS" || v == "-SH" {
			return []Violation{{
				KataID: "ZC1418",
				Message: "Use Zsh `limit -h` (hard) / `limit -s` (soft) instead of " +
					"`ulimit -H`/`-S`. Zsh's `limit` builtin is more human-readable.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
