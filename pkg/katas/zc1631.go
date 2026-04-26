// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1631",
		Title:    "Error on `openssl ... -passin pass:SECRET` / `-passout pass:SECRET`",
		Severity: SeverityError,
		Description: "OpenSSL's `-passin` / `-passout` accept a password source selector. The " +
			"`pass:LITERAL` form embeds the password as an argv element — visible in `ps`, " +
			"`/proc/<pid>/cmdline`, shell history, and audit logs. Use one of the safer " +
			"sources: `env:VARNAME` reads from an env var, `file:PATH` reads the first line " +
			"of PATH, `fd:N` reads from an open descriptor, `stdin` reads a line from stdin.",
		Check: checkZC1631,
	})
}

func checkZC1631(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "openssl" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v != "-passin" && v != "-passout" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		val := cmd.Arguments[i+1].String()
		if !strings.HasPrefix(val, "pass:") {
			continue
		}
		return []Violation{{
			KataID: "ZC1631",
			Message: "`openssl " + v + " " + val + "` puts the password in argv — visible " +
				"via `ps`. Use `env:VARNAME`, `file:PATH`, `fd:N`, or `stdin`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}
