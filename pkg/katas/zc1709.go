// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1709",
		Title:    "Error on `htpasswd -b USER PASSWORD` — basic-auth password in process list",
		Severity: SeverityError,
		Description: "`htpasswd -b FILE USER PASSWORD` (batch mode) takes the password as an " +
			"argv slot. The cleartext sits in `/proc/PID/cmdline`, shell history, audit " +
			"records, and any `ps` output. Use `htpasswd -i FILE USER` and pipe the " +
			"secret on stdin (`printf %s \"$pw\" | htpasswd -i FILE USER`), or omit `-b` " +
			"and `-i` so htpasswd prompts on the controlling TTY.",
		Check: checkZC1709,
	})
}

func checkZC1709(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "htpasswd" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !strings.HasPrefix(v, "-") || strings.HasPrefix(v, "--") {
			continue
		}
		body := strings.TrimPrefix(v, "-")
		if strings.ContainsRune(body, 'b') {
			return []Violation{{
				KataID: "ZC1709",
				Message: "`htpasswd -b USER PASSWORD` puts the password in argv — visible " +
					"via `ps` / `/proc/PID/cmdline`. Use `htpasswd -i FILE USER` with the " +
					"password piped on stdin instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
