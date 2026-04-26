// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1490",
		Title:    "Error on `socat ... EXEC:<shell>` / `SYSTEM:<shell>` — socat reverse-shell pattern",
		Severity: SeverityError,
		Description: "The `EXEC:` and `SYSTEM:` socat address types spawn a subprocess connected " +
			"to the other socat endpoint. Paired with `TCP:` or `TCP-LISTEN:`, they form the " +
			"second-most-common reverse/bind shell payload after `nc -e`. Legitimate uses exist " +
			"(test harnesses, pty brokers) but should be gated behind explicit authorization " +
			"and a non-shell command. Scan hits are worth a look.",
		Check: checkZC1490,
	})
}

func checkZC1490(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "socat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		low := strings.ToLower(v)
		// EXEC:/bin/bash or "EXEC:\"/bin/sh -i\",pty,stderr"
		if strings.Contains(low, "exec:/bin/bash") ||
			strings.Contains(low, "exec:/bin/sh") ||
			strings.Contains(low, "exec:/bin/zsh") ||
			strings.Contains(low, "exec:\"/bin/bash") ||
			strings.Contains(low, "exec:\"/bin/sh") ||
			strings.Contains(low, "system:/bin/bash") ||
			strings.Contains(low, "system:/bin/sh") {
			return []Violation{{
				KataID: "ZC1490",
				Message: "`socat` pointed at a shell via `EXEC:` / `SYSTEM:` — matches the " +
					"classic reverse/bind-shell pattern. Gate behind explicit authorization.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
