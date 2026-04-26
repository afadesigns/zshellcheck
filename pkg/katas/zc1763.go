// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1763",
		Title:    "Error on `docker compose down -v` / `--volumes` — wipes named volumes (data loss)",
		Severity: SeverityError,
		Description: "`docker compose down -v` (alias `--volumes`, equivalent in `docker-compose " +
			"down -v`) tears the stack down AND deletes every named volume declared in the " +
			"compose file. Database contents, cache state, uploaded assets, and any other " +
			"volume-backed data goes with them — there is no soft-delete. Drop the flag in " +
			"CI and production scripts; keep it only for throwaway local testbeds where " +
			"losing volume state is intentional.",
		Check: checkZC1763,
	})
}

func checkZC1763(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var argsAfterDown []ast.Expression
	switch ident.Value {
	case "docker":
		if len(cmd.Arguments) < 2 {
			return nil
		}
		if cmd.Arguments[0].String() != "compose" || cmd.Arguments[1].String() != "down" {
			return nil
		}
		argsAfterDown = cmd.Arguments[2:]
	case "docker-compose":
		if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "down" {
			return nil
		}
		argsAfterDown = cmd.Arguments[1:]
	default:
		return nil
	}

	for _, arg := range argsAfterDown {
		v := arg.String()
		if v == "-v" || v == "--volumes" {
			return []Violation{{
				KataID: "ZC1763",
				Message: "`docker compose down " + v + "` wipes every named volume declared " +
					"in the stack — database, cache, uploaded assets go with it. Drop " +
					"the flag in CI / prod scripts.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
