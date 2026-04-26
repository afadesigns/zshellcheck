// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1622BashOps = []string{"@U}", "@L}", "@Q}", "@E}", "@A}", "@K}", "@a}"}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1622",
		Title:    "Style: `${var@U/L/Q/...}` — prefer Zsh `${(U)var}` / `${(L)var}` / `${(Q)var}` flags",
		Severity: SeverityStyle,
		Description: "The `@<op>` suffix came from Bash 5. Zsh 5.9+ compiles in compatibility " +
			"for the common ones, but the idiomatic Zsh form is the `(X)var` parameter-" +
			"expansion flag — `${(U)var}` uppercase, `${(L)var}` lowercase, `${(Q)var}` " +
			"unquote, `${(k)var}` keys, `${(t)var}` type, `${(e)var}` re-evaluate. The flag " +
			"form composes (`${(Uf)str}` works) and reads consistently across the Zsh " +
			"documentation. Prefer the native flag over the Bash-compat form.",
		Check: checkZC1622,
	})
}

func checkZC1622(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !strings.Contains(v, "${") {
			continue
		}
		for _, op := range zc1622BashOps {
			if strings.Contains(v, op) {
				return []Violation{{
					KataID: "ZC1622",
					Message: "`${var" + strings.TrimSuffix(op, "}") + "}` — prefer Zsh " +
						"`${(X)var}` parameter-expansion flags (e.g. `${(U)var}` for " +
						"uppercase).",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityStyle,
				}}
			}
		}
	}
	return nil
}
