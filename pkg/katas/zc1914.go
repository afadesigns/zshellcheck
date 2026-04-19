package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1914",
		Title:    "Warn on `curl --doh-url …` / `--dns-servers …` — overrides system resolver per-request",
		Severity: SeverityWarning,
		Description: "`curl --doh-url https://doh.example/dns-query` routes the lookup through a " +
			"caller-specified DNS-over-HTTPS endpoint; `curl --dns-servers 1.1.1.1,8.8.8.8` " +
			"forces classic UDP to the listed servers. Both detour around the host's resolver " +
			"chain — `/etc/hosts`, `systemd-resolved`, `nsswitch`, split-horizon DNS — so the " +
			"request lands at an IP the operator did not vet. In production scripts that is " +
			"usually a mis-debug line left in; drop the flag or gate it behind an explicit " +
			"`--doh-insecure` + `--resolve` pinning audit so reviewers can see the intent.",
		Check: checkZC1914,
	})
}

func checkZC1914(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `curl --doh-url https://…` may mangle the name.
	switch ident.Value {
	case "doh-url":
		return zc1914Hit(cmd, "--doh-url")
	case "dns-servers":
		return zc1914Hit(cmd, "--dns-servers")
	case "curl":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			switch {
			case v == "--doh-url", strings.HasPrefix(v, "--doh-url="):
				return zc1914Hit(cmd, "--doh-url")
			case v == "--dns-servers", strings.HasPrefix(v, "--dns-servers="):
				return zc1914Hit(cmd, "--dns-servers")
			}
		}
	}
	return nil
}

func zc1914Hit(cmd *ast.SimpleCommand, flag string) []Violation {
	return []Violation{{
		KataID: "ZC1914",
		Message: "`curl " + flag + "` bypasses the host's resolver chain — `/etc/hosts`, " +
			"`systemd-resolved`, split-horizon DNS — so the request lands at an IP the " +
			"operator did not vet. Drop the flag or pair it with `--resolve` pinning.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
