package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1382",
		Title:    "Avoid `$READLINE_LINE`/`$READLINE_POINT` — Zsh ZLE uses `$BUFFER`/`$CURSOR`",
		Severity: SeverityError,
		Description: "Bash readline exposes the current input line as `$READLINE_LINE` and cursor " +
			"offset as `$READLINE_POINT` inside `bind -x` handlers. Zsh's Line Editor (ZLE) uses " +
			"`$BUFFER` (line text) and `$CURSOR` (1-based column) inside widget functions. The " +
			"Bash names are unset in Zsh.",
		Check: checkZC1382,
	})
}

func checkZC1382(node ast.Node) []Violation {
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
		if strings.Contains(v, "READLINE_LINE") || strings.Contains(v, "READLINE_POINT") ||
			strings.Contains(v, "READLINE_MARK") {
			return []Violation{{
				KataID: "ZC1382",
				Message: "Bash `$READLINE_*` vars do not exist in Zsh. Inside ZLE widgets use " +
					"`$BUFFER`, `$CURSOR`, `$MARK`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
