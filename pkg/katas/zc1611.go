package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1611",
		Title:    "Style: `${var^^}` / `${var,,}` — prefer Zsh `${(U)var}` / `${(L)var}` for case change",
		Severity: SeverityStyle,
		Description: "`${var^^}` (uppercase) and `${var,,}` (lowercase) came from Bash 4. Zsh " +
			"accepts them for compatibility but the idiomatic form is the parameter-expansion " +
			"flag: `${(U)var}` / `${(L)var}`. The flag is also available per-element in " +
			"arrays (`${(U)array}`) and composes with other flags (`${(UL)array}` doesn't " +
			"make sense, but `${(U)${(f)str}}` does). Prefer the Zsh-native form in a `.zsh` " +
			"script; it keeps the codebase consistent with other `(X)var` patterns.",
		Check: checkZC1611,
	})
}

func checkZC1611(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !strings.Contains(v, "${") {
			continue
		}
		if strings.Contains(v, "^^}") || strings.Contains(v, ",,}") {
			return []Violation{{
				KataID: "ZC1611",
				Message: "`${var^^}` / `${var,,}` — prefer Zsh `${(U)var}` / `${(L)var}` " +
					"for case conversion.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}
	return nil
}
