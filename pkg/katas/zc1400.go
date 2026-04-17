package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1400",
		Title:    "Use Zsh `$CPUTYPE` for architecture detection instead of parsing `$HOSTTYPE`",
		Severity: SeverityInfo,
		Description: "Bash's `$HOSTTYPE` is a combined architecture/vendor/OS string (e.g. " +
			"`x86_64-pc-linux-gnu`). Zsh exposes the same as `$HOSTTYPE` but additionally splits " +
			"out `$CPUTYPE` (e.g. `x86_64`) for pure architecture queries — no `awk -F-` " +
			"needed to extract.",
		Check: checkZC1400,
	})
}

func checkZC1400(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	// Only fire when HOSTTYPE is being parsed/split (cut, awk, sed usage).
	switch ident.Value {
	case "cut", "awk", "sed":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "HOSTTYPE") {
			return []Violation{{
				KataID: "ZC1400",
				Message: "Use Zsh `$CPUTYPE` for pure architecture instead of splitting " +
					"`$HOSTTYPE` with `cut`/`awk`/`sed`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}
