// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1673",
		Title:    "Style: `stty -echo` around `read` — prefer Zsh `read -s`",
		Severity: SeverityStyle,
		Description: "The classic `stty -echo; IFS= read -r password; stty echo` pattern has a " +
			"serious failure mode: a crash or SIGINT between the two `stty` calls leaves " +
			"the user's terminal stuck in echo-off, which is silent and confusing. Zsh's " +
			"`read -s VAR` (also Bash 4+) disables echo only for that one `read`, restores " +
			"it on return even if the read is interrupted, and avoids two external forks. " +
			"Switch the prompt to `read -s` (or `read -ks` for single-key password) and " +
			"drop the `stty` bracketing.",
		Check: checkZC1673,
	})
}

func checkZC1673(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "stty" {
		return nil
	}

	if len(cmd.Arguments) != 1 {
		return nil
	}
	if cmd.Arguments[0].String() != "-echo" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1673",
		Message: "`stty -echo` to mask password entry is fragile — a crash leaves the " +
			"terminal echo-off. Use `read -s VAR` (Zsh / Bash 4+) instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
