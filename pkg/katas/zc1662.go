package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1662",
		Title:    "Error on `pkexec env VAR=VAL CMD` — controlled env crossed into the root session",
		Severity: SeverityError,
		Description: "`pkexec env VAR=VALUE CMD` invokes `/usr/bin/env` as the target user (root " +
			"by default) with a caller-controlled environment. Polkit sanitizes a short " +
			"allow-list on its own, but once `env` takes over the remaining variables " +
			"(`LD_PRELOAD`, `GCONV_PATH`, `PYTHONPATH`, `XDG_RUNTIME_DIR`, `LANGUAGE`) ride " +
			"straight into root. CVE-2021-4034 (pwnkit) demonstrated the same primitive by " +
			"abusing argv[0]; the `env` wrapper makes the bypass trivial. If the child needs " +
			"specific variables, set them in a polkit rule or via `systemd-run --user` " +
			"instead, not through `env`.",
		Check: checkZC1662,
	})
}

func checkZC1662(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pkexec" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "env" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1662",
		Message: "`pkexec env VAR=VAL CMD` hands the root session a caller-controlled " +
			"environment — use a polkit rule or `systemd-run --user` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
