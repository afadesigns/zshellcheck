// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1846",
		Title:    "Warn on `certbot … --force-renewal` — bypasses ACME rate-limit safety",
		Severity: SeverityWarning,
		Description: "`certbot renew --force-renewal` and `certbot certonly --force-renewal` reissue " +
			"a certificate regardless of remaining validity. Placed in a daily cron, the " +
			"same hostname burns through Let's Encrypt's per-domain rate limits (50 " +
			"certificates per registered domain per 7 days, 5 duplicate certificates per " +
			"domain per 7 days); once the limit trips, no cert for that host — fresh or " +
			"renewal — can be issued until the rolling window expires, which often happens " +
			"during an outage when you need it least. Drop `--force-renewal` and let " +
			"certbot's default 30-days-before-expiry gate do its job, or if you really need " +
			"a specific reissue, run it once manually.",
		Check: checkZC1846,
	})
}

func checkZC1846(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "certbot" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}
	sub := args[0].String()
	if sub != "renew" && sub != "certonly" && sub != "run" {
		return nil
	}
	for _, arg := range args[1:] {
		if arg.String() == "--force-renewal" {
			return []Violation{{
				KataID: "ZC1846",
				Message: "`certbot " + sub + " --force-renewal` reissues regardless of " +
					"expiry — in a cron it burns Let's Encrypt rate limits (50 certs " +
					"per domain / 7 days). Drop the flag and let the 30-day gate work.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
