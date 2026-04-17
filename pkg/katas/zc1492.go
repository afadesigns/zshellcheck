package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1492",
		Title:    "Style: `at` / `batch` for deferred execution — prefer systemd timers for auditability",
		Severity: SeverityStyle,
		Description: "`at` and `batch` schedule one-shot deferred jobs via `atd`. The job payload " +
			"lands in `/var/spool/at*/` with no unit file or dependency graph, which makes it " +
			"harder to review in fleet audits, easier to miss in a compromise triage, and one of " +
			"the less-watched places adversaries stash persistence. Prefer `systemd-run " +
			"--on-calendar=` or a proper `.timer` unit with a corresponding `.service`.",
		Check: checkZC1492,
	})
}

func checkZC1492(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "at" {
		return nil
	}

	// Skip list/remove/query forms: -l, -r, -c, -q
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-l" || v == "-r" || v == "-c" || v == "-q" ||
			v == "--list" || v == "--remove" {
			return nil
		}
	}

	// Any remaining `at <time>` or `at -f <script> <time>` invocation.
	if len(cmd.Arguments) == 0 {
		return nil
	}
	return []Violation{{
		KataID: "ZC1492",
		Message: "`at` schedules via atd with no unit file — harder to audit. Prefer " +
			"`systemd-run --on-calendar=` or a `.timer` unit.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
