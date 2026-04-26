// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1553",
		Title:    "Style: use Zsh `${(U)var}` / `${(L)var}` instead of `tr '[:lower:]' '[:upper:]'`",
		Severity: SeverityStyle,
		Description: "Zsh provides `${(U)var}` and `${(L)var}` parameter-expansion flags for " +
			"case conversion in-process. Spawning `tr` for this forks/execs per call (noticeable " +
			"in a hot loop), relies on the external `tr` being POSIX-compliant (BusyBox and old " +
			"macOS differ), and round-trips the data through a pipe. Drop `tr` for the " +
			"built-in: `upper=${(U)lower}` / `lower=${(L)upper}`.",
		Check: checkZC1553,
	})
}

func checkZC1553(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "tr" {
		return nil
	}
	from, to := zc1553TrSets(cmd)
	if !zc1553IsCasePair(from, to) {
		return nil
	}
	return []Violation{{
		KataID: "ZC1553",
		Message: "`tr` for case conversion — use Zsh `${(U)var}` / `${(L)var}` to avoid " +
			"the fork/exec and portability hazard.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

func zc1553TrSets(cmd *ast.SimpleCommand) (from, to string) {
	for _, a := range cmd.Arguments {
		v := strings.Trim(a.String(), "'\"")
		if strings.HasPrefix(v, "-") {
			continue
		}
		if from == "" {
			from = v
			continue
		}
		if to == "" {
			to = v
			return
		}
	}
	return
}

func zc1553IsCasePair(from, to string) bool {
	upper := from == "[:upper:]" || from == "A-Z"
	lower := from == "[:lower:]" || from == "a-z"
	if !upper && !lower {
		return false
	}
	other := to == "[:upper:]" || to == "A-Z" || to == "[:lower:]" || to == "a-z"
	return other && upper != (to == "[:upper:]" || to == "A-Z")
}
