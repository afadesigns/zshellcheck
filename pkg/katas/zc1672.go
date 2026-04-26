// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1672",
		Title:    "Info: `chcon` writes an ephemeral SELinux label — next `restorecon` wipes it",
		Severity: SeverityInfo,
		Description: "`chcon -t TYPE PATH` sets the file context out-of-band; it does not update " +
			"the `file_contexts` policy database. As soon as `restorecon`, `semodule -n`, or " +
			"a policy rebuild runs, the label snaps back to whatever the compiled policy " +
			"says — often `default_t`, which can break a deployed workload or silently " +
			"re-introduce a denial the script tried to fix. For anything long-lived use " +
			"`semanage fcontext -a -t TYPE '<regex>'` then `restorecon -F <path>` so the " +
			"mapping lives in policy.",
		Check: checkZC1672,
	})
}

func checkZC1672(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chcon" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1672",
		Message: "`chcon` writes an ephemeral SELinux label — `restorecon` / policy rebuild " +
			"reverts it. Persist via `semanage fcontext -a -t TYPE 'REGEX'` + `restorecon`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}
