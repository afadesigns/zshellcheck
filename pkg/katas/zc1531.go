package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1531",
		Title:    "Warn on `wget -t 0` — infinite retries, hangs on a dead endpoint",
		Severity: SeverityWarning,
		Description: "`wget -t 0` (or `--tries=0`) means retry forever. Paired with `-w` (wait " +
			"between retries) and a dead endpoint, the script hangs until killed — in a cron " +
			"job, every subsequent invocation piles up and eventually the UID's process limit " +
			"trips. Use a finite retry count (`-t 5`) plus `--timeout=<seconds>` to cap total " +
			"wall time.",
		Check: checkZC1531,
	})
}

func checkZC1531(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "wget" {
		return nil
	}

	var prevT bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevT {
			prevT = false
			if v == "0" {
				return []Violation{{
					KataID: "ZC1531",
					Message: "`wget -t 0` retries forever — script hangs on dead endpoint. " +
						"Use finite `-t 5` plus `--timeout=<seconds>`.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-t" {
			prevT = true
		}
	}
	return nil
}
