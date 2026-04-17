package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1381",
		Title:    "Avoid `$COMP_WORDS`/`$COMP_CWORD` — Zsh uses `words`/`$CURRENT`",
		Severity: SeverityError,
		Description: "Bash programmable completion reads the partial command via `$COMP_WORDS` " +
			"(array of tokens) and `$COMP_CWORD` (index of cursor). Zsh's completion system " +
			"exposes the same via `words` (array) and `$CURRENT` (1-based cursor index). Using " +
			"the Bash names in Zsh completion functions produces empty expansions.",
		Check: checkZC1381,
	})
}

func checkZC1381(node ast.Node) []Violation {
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
		if strings.Contains(v, "COMP_WORDS") || strings.Contains(v, "COMP_CWORD") ||
			strings.Contains(v, "COMP_LINE") || strings.Contains(v, "COMP_POINT") {
			return []Violation{{
				KataID: "ZC1381",
				Message: "Bash `$COMP_*` completion variables do not exist in Zsh. Use " +
					"`$words` (array of tokens), `$CURRENT` (cursor index), `$BUFFER`, or the " +
					"`_arguments`/`_values` helpers from `compsys`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
