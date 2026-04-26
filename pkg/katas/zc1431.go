// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1431",
		Title:    "Dangerous: `crontab -r` — removes all the user's cron jobs without confirmation",
		Severity: SeverityWarning,
		Description: "`crontab -r` deletes the entire crontab for the current user (or the target " +
			"user with `-u`). There is no `.bak` left behind, no `-i` prompt by default on most " +
			"platforms. Back up first with `crontab -l > /tmp/cron.bak`, then use `crontab -ir` " +
			"(interactive) to require confirmation.",
		Check: checkZC1431,
	})
}

func checkZC1431(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "crontab" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-r" || v == "-ur" || v == "-ru" {
			return []Violation{{
				KataID: "ZC1431",
				Message: "`crontab -r` removes all cron jobs with no backup. Save first " +
					"(`crontab -l > cron.bak`) and use `crontab -ir` for interactive confirmation.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
