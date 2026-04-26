// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1718",
		Title:    "Error on `gh secret set --body SECRET` / `-b SECRET` — secret in process list",
		Severity: SeverityError,
		Description: "`gh secret set NAME --body VALUE` (or `-b VALUE`, `--body=VALUE`) puts the " +
			"secret on the command line. The cleartext appears in `ps`, `/proc/<pid>/cmdline`, " +
			"shell history, and the audit log of the host running `gh`. Pipe the value via " +
			"stdin (`gh secret set NAME < file`, `printf %s \"$SECRET\" | gh secret set NAME " +
			"--body -`) or use `--body-file PATH` so the value never lands in argv.",
		Check: checkZC1718,
	})
}

func checkZC1718(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "secret" || cmd.Arguments[1].String() != "set" {
		return nil
	}

	prevBody := false
	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if prevBody {
			if v == "-" {
				return nil
			}
			return zc1718Hit(cmd, "--body "+v)
		}
		switch {
		case v == "--body" || v == "-b":
			prevBody = true
		case strings.HasPrefix(v, "--body="):
			val := strings.TrimPrefix(v, "--body=")
			if val == "-" {
				return nil
			}
			return zc1718Hit(cmd, v)
		}
	}
	return nil
}

func zc1718Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1718",
		Message: "`gh secret set ... " + what + "` puts the secret in argv — visible in " +
			"`ps`, `/proc`, history. Use `--body-file PATH` or pipe via stdin " +
			"(`... --body -` with `printf %s \"$SECRET\" |`).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
