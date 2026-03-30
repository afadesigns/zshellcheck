package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1152",
		Title:    "Use Zsh PCRE module instead of `grep -P`",
		Severity: SeverityStyle,
		Description: "`grep -P` (Perl regex) is not available on all platforms (e.g., macOS). " +
			"Use `zmodload zsh/pcre` and `pcre_compile`/`pcre_match` for portable PCRE matching.",
		Check: checkZC1152,
	})
}

func checkZC1152(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-P" {
			return []Violation{{
				KataID: "ZC1152",
				Message: "Avoid `grep -P` — it's unavailable on macOS. Use `zmodload zsh/pcre` " +
					"with `pcre_compile`/`pcre_match` or `grep -E` for portable regex matching.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
