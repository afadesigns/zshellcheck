package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1579",
		Title:    "Warn on `curl --retry-all-errors` without `--max-time` — hammers endpoint on failure",
		Severity: SeverityWarning,
		Description: "`--retry-all-errors` (curl 7.71+) treats every HTTP error as retryable. " +
			"Without `--max-time` capping total wall clock, a server that responds `500` " +
			"quickly gets hit back-to-back until `--retry` exhausts — a mini-DoS against your " +
			"own upstream, especially if the script itself is scheduled on many nodes. Pair " +
			"with `--max-time <seconds>` or prefer `--retry-connrefused` (only retries " +
			"connection-level failures).",
		Check: checkZC1579,
	})
}

func checkZC1579(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "curl" {
		return nil
	}

	var hasRetryAll, hasMaxTime bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--retry-all-errors" {
			hasRetryAll = true
		}
		if v == "--max-time" || v == "-m" {
			hasMaxTime = true
		}
	}
	if !hasRetryAll || hasMaxTime {
		return nil
	}
	return []Violation{{
		KataID: "ZC1579",
		Message: "`curl --retry-all-errors` with no `--max-time` hammers the upstream on " +
			"failure. Pair with `-m <seconds>` or use `--retry-connrefused`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
