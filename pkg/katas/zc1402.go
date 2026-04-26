// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1402",
		Title:    "Avoid `date -d @seconds` — use Zsh `strftime` for epoch formatting",
		Severity: SeverityStyle,
		Description: "`date -d @N -- '+fmt'` / `date --date=@N` converts epoch seconds to a " +
			"formatted date. Zsh's `zsh/datetime` module provides `strftime fmt N` directly " +
			"— a single builtin, no `date` spawn, and the `-d`/`@` form is GNU-specific " +
			"(not portable to BSD `date`).",
		Check: checkZC1402,
	})
}

func checkZC1402(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "date" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-d" || v == "--date" ||
			(len(v) > 6 && v[:6] == "--date=") {
			return []Violation{{
				KataID: "ZC1402",
				Message: "Use Zsh `strftime` (from `zsh/datetime`) instead of `date -d @N -- +fmt`. " +
					"The `-d`/`@` form is GNU-specific; `strftime` is portable Zsh.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
