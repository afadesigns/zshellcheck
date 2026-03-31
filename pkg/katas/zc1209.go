package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1209",
		Title:    "Use `systemctl --no-pager` in scripts",
		Severity: SeverityStyle,
		Description: "`systemctl` invokes a pager by default which hangs in non-interactive scripts. " +
			"Use `--no-pager` or pipe to `cat` for reliable script output.",
		Check: checkZC1209,
	})
}

func checkZC1209(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "--no-pager" {
			return nil
		}
	}

	// Only flag subcommands that produce output (status, list-units, etc.)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "status" || val == "list-units" || val == "list-timers" || val == "show" {
			return []Violation{{
				KataID: "ZC1209",
				Message: "Use `systemctl --no-pager` in scripts. Without it, " +
					"systemctl invokes a pager that hangs in non-interactive execution.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
