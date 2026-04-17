package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1401",
		Title:    "Prefer Zsh `$VENDOR` over parsing `$MACHTYPE` for vendor detection",
		Severity: SeverityInfo,
		Description: "Both Bash and Zsh expose `$MACHTYPE` (e.g. `x86_64-pc-linux-gnu`). Zsh " +
			"additionally pre-parses the vendor component into `$VENDOR` (e.g. `pc`, `apple`). " +
			"Avoid `cut -d- -f2 <<< $MACHTYPE` when `$VENDOR` is available directly.",
		Check: checkZC1401,
	})
}

func checkZC1401(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "cut", "awk", "sed":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "MACHTYPE") {
			return []Violation{{
				KataID: "ZC1401",
				Message: "Use Zsh `$VENDOR` for vendor field instead of splitting `$MACHTYPE` " +
					"with `cut`/`awk`/`sed`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}
