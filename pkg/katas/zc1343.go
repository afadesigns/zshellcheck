package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1343",
		Title:    "Use Zsh `*(m±N)` glob qualifier instead of `find -mtime N`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(mN)`, `*(m+N)`, `*(m-N)` glob qualifiers match files by age in days " +
			"(exact / older / newer). For hours use `*(h±N)`, for minutes `*(M±N)`. " +
			"Same expressive power as `find -mtime`, no external process.",
		Check: checkZC1343,
	})
}

func checkZC1343(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-mtime" || val == "-mmin" || val == "-atime" || val == "-amin" ||
			val == "-ctime" || val == "-cmin" {
			return []Violation{{
				KataID: "ZC1343",
				Message: "Use Zsh glob qualifiers (`*(m±N)`, `*(M±N)`, `*(a±N)`, `*(c±N)`) instead of " +
					"`find -mtime`/`-mmin`/`-atime`/`-amin`/`-ctime`/`-cmin` for age predicates.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
