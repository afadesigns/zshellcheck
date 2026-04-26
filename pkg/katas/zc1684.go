// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1684",
		Title:    "Error on `redis-cli -a PASSWORD` — authentication password in process list",
		Severity: SeverityError,
		Description: "`redis-cli -a <password>` (and the joined form `-aPASSWORD`) puts the " +
			"authentication password in the command line — visible to every user on the " +
			"host through `ps`, `/proc/PID/cmdline`, audit logs, and shell history. redis-" +
			"cli 6.0+ prints a warning to stderr but still connects. Use the " +
			"`REDISCLI_AUTH` environment variable (read automatically by redis-cli), or " +
			"`-askpass` to prompt from TTY; both keep the secret out of the argv tail.",
		Check: checkZC1684,
	})
}

func checkZC1684(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "redis-cli" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-a" && i+1 < len(cmd.Arguments) {
			return zc1684Hit(cmd)
		}
		if strings.HasPrefix(v, "-a") && v != "-a" && !strings.HasPrefix(v, "--") {
			return zc1684Hit(cmd)
		}
	}
	return nil
}

func zc1684Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1684",
		Message: "`redis-cli -a PASSWORD` leaks the password into `ps` / `/proc/PID/cmdline` — " +
			"use `REDISCLI_AUTH` env var or `-askpass` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
