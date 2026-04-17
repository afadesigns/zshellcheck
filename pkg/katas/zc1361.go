package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1361",
		Title:    "Avoid `awk 'NR==N'` — use Zsh array subscript on `${(f)...}`",
		Severity: SeverityStyle,
		Description: "Picking the N-th line with `awk 'NR==N'` spawns awk. Zsh can split file " +
			"contents on newlines with `${(f)\"$(<file)\"}` and index directly: `lines=(${(f)\"$(<f)\"}); print $lines[N]`.",
		Check: checkZC1361,
	})
}

func checkZC1361(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "awk" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		// Look for NR== or NR == in the awk program
		if strings.Contains(val, "NR==") || strings.Contains(val, "NR ==") {
			return []Violation{{
				KataID: "ZC1361",
				Message: "Avoid `awk 'NR==N'` — split with `${(f)\"$(<file)\"}` in Zsh and index: " +
					"`lines=(${(f)\"$(<file)\"}); print $lines[N]`. No awk process needed.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
