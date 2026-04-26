// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1574",
		Title:    "Warn on `git config credential.helper store` — plaintext credentials on disk",
		Severity: SeverityWarning,
		Description: "`credential.helper store` writes the username and password to " +
			"`~/.git-credentials` in plaintext. Anything that backs up that file (rsync, " +
			"imaging, cloud sync) then carries the credential around. Use a platform helper " +
			"instead: `manager` / `manager-core` on Windows / Mac, `libsecret` on Linux, or " +
			"`cache --timeout=3600` for short-lived in-memory caching.",
		Check: checkZC1574,
	})
}

func checkZC1574(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "git" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// git config [scope] credential.helper store
	for i, a := range args {
		if a != "config" {
			continue
		}
		// Walk past optional scope flag.
		j := i + 1
		for j < len(args) && strings.HasPrefix(args[j], "--") && args[j] != "--" {
			j++
		}
		if j+1 < len(args) && args[j] == "credential.helper" {
			v := args[j+1]
			// Accept bare `store` or quoted `store --file=...`.
			if v == "store" || strings.HasPrefix(v, "store ") ||
				strings.Contains(v, "store --file=") ||
				strings.Contains(v, "'store") || strings.Contains(v, "\"store") {
				return []Violation{{
					KataID: "ZC1574",
					Message: "`git credential.helper store` saves credentials in plaintext — " +
						"backups leak the token. Use platform helper (manager-core / " +
						"libsecret) or `cache --timeout=<sec>`.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}
