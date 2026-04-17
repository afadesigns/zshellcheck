package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1487",
		Title:    "Warn on `history -c` — clears shell history (and is a Bash-ism under Zsh)",
		Severity: SeverityWarning,
		Description: "`history -c` clears the in-memory history buffer in Bash. It is a standard " +
			"post-compromise anti-forensics step. It is also a Bash-ism: in Zsh, `history` " +
			"takes completely different arguments, so a copy-pasted `history -c` silently no-ops " +
			"and leaves the author thinking history was cleared when it was not. If you really " +
			"need to rotate history in a Zsh script, unset `HISTFILE` before the sensitive " +
			"block or redirect to `/dev/null` explicitly.",
		Check: checkZC1487,
	})
}

func checkZC1487(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "history" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" || v == "-d" {
			return []Violation{{
				KataID: "ZC1487",
				Message: "`history " + v + "` is a Bash-ism for clearing history — does " +
					"nothing in Zsh and is a classic post-compromise tactic elsewhere.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
