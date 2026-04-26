// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1577",
		Title:    "Warn on `dig <name> ANY` — deprecated query type (RFC 8482)",
		Severity: SeverityWarning,
		Description: "ANY queries return whatever the authoritative server feels like sending " +
			"back — or just the HINFO placeholder mandated by RFC 8482. Modern recursors " +
			"filter ANY to avoid reflection-amplification abuse, so scripts that rely on ANY " +
			"for enumeration get inconsistent or empty results. Query the specific record " +
			"types you want (`dig A name`, `dig MX name`, `dig NS name`) and combine them.",
		Check: checkZC1577,
	})
}

func checkZC1577(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dig" && ident.Value != "drill" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "ANY" || v == "any" {
			return []Violation{{
				KataID: "ZC1577",
				Message: "`" + ident.Value + " ... ANY` is RFC 8482-deprecated — filtered by " +
					"recursors. Query specific types (A / MX / NS / …) and combine.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
