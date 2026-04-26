// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1775",
		Title:    "Warn on `timeout DURATION cmd` without `--kill-after` / `-k` — hang on SIGTERM-resistant child",
		Severity: SeverityWarning,
		Description: "`timeout DURATION cmd` sends `SIGTERM` once the duration elapses and then " +
			"waits for the child to exit. A child that blocks or ignores `SIGTERM` (long-running " +
			"daemons, processes stuck in `D` state, a trapped / reset signal handler) never " +
			"dies, so the entire pipeline hangs past the intended bound. Add `--kill-after=N` " +
			"(`-k N`) so timeout escalates to `SIGKILL` after N seconds, guaranteeing exit. " +
			"Typical choice: a few seconds shorter than your CI step budget, so the overall " +
			"wait remains bounded.",
		Check: checkZC1775,
	})
}

func checkZC1775(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "timeout" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-k" || v == "--kill-after" {
			return nil
		}
		if strings.HasPrefix(v, "--kill-after=") {
			return nil
		}
	}
	return []Violation{{
		KataID: "ZC1775",
		Message: "`timeout` without `--kill-after` / `-k` only sends `SIGTERM` — a child " +
			"that blocks or ignores it hangs the pipeline past the deadline. Add " +
			"`--kill-after=N` so timeout escalates to `SIGKILL`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
