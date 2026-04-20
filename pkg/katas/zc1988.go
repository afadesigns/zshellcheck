package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1988",
		Title:    "Error on `nsupdate -y HMAC:NAME:SECRET` — TSIG key visible in argv and shell history",
		Severity: SeverityError,
		Description: "`nsupdate -y [alg:]name:base64secret` hands the TSIG shared secret " +
			"directly on the command line, so `ps auxf`, `/proc/PID/cmdline`, and " +
			"`$HISTFILE` all capture the key — and whoever owns the key can rewrite " +
			"any zone that trusts it (DNS hijack, MX hijack, ACME domain-validation " +
			"bypass). `nsupdate -k /etc/named/KEY` (or `-k $KEYFILE` with `0600` " +
			"perms) reads the key from disk without exposing it. If the secret must " +
			"come from a secret store, pipe it through `nsupdate -k /dev/stdin <<<\"$KEYFILE_CONTENTS\"` " +
			"so the raw material never lands in argv.",
		Check: checkZC1988,
	})
}

func checkZC1988(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "nsupdate" {
		return nil
	}
	for i, arg := range cmd.Arguments {
		if arg.String() == "-y" && i+1 < len(cmd.Arguments) {
			return []Violation{{
				KataID: "ZC1988",
				Message: "`nsupdate -y …` puts the TSIG key in argv — `ps`, " +
					"`/proc/*/cmdline`, and `$HISTFILE` all capture it. Use " +
					"`nsupdate -k $KEYFILE` with a `0600` keyfile instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
