// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1722",
		Title:    "Warn on `ssh-keyscan HOST >> known_hosts` — TOFU bypass, blind-trust new host key",
		Severity: SeverityWarning,
		Description: "`ssh-keyscan` fetches whatever host key the remote serves on its first reply. " +
			"Appending the result straight to `known_hosts` is the exact step the host-key " +
			"check is meant to defend against: a man-in-the-middle on first contact wins " +
			"permanently. Pin the expected fingerprint via a side channel (vendor docs, prior " +
			"verified contact) and assert it matches `ssh-keyscan HOST | ssh-keygen -lf -` " +
			"before the append.",
		Check: checkZC1722,
	})
}

func checkZC1722(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ssh-keyscan" {
		return nil
	}

	prevRedir := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevRedir != "" {
			if strings.Contains(v, "known_hosts") {
				return []Violation{{
					KataID: "ZC1722",
					Message: "`ssh-keyscan ... " + prevRedir + " " + v + "` accepts the " +
						"first-served host key without verifying its fingerprint. Pipe " +
						"to `ssh-keygen -lf -` and assert the fingerprint first.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
			prevRedir = ""
			continue
		}
		if v == ">>" || v == ">" {
			prevRedir = v
		}
	}
	return nil
}
