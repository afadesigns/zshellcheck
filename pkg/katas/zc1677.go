package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1677",
		Title:    "Warn on `trap 'set -x' DEBUG` — xtrace on every command leaks secrets",
		Severity: SeverityWarning,
		Description: "`trap 'set -x' DEBUG` runs the trap handler before every simple command, " +
			"turning on xtrace for the remainder of the shell. Every subsequent `curl " +
			"-H 'Authorization: Bearer …'`, `mysql -p<password>`, or `aws configure set " +
			"…` then prints its full argv to stderr — commonly into a log file or CI " +
			"artifact. The same antipattern shows up as `set -o xtrace` inside a DEBUG " +
			"trap. Instrument selectively with `typeset -ft FUNC` (Zsh function-level " +
			"xtrace), or add `exec 2>>\"$log\"; set -x` only around the part of the " +
			"script you want traced.",
		Check: checkZC1677,
	})
}

func checkZC1677(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "trap" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	last := cmd.Arguments[len(cmd.Arguments)-1].String()
	if last != "DEBUG" {
		return nil
	}

	handler := strings.Trim(cmd.Arguments[0].String(), "'\"")
	if !strings.Contains(handler, "set -x") && !strings.Contains(handler, "set -o xtrace") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1677",
		Message: "`trap 'set -x' DEBUG` keeps xtrace on after the first command — " +
			"every subsequent argv (passwords, bearer tokens) lands in the log. Trace a " +
			"narrow block with `set -x … set +x` or use `typeset -ft FUNC` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
