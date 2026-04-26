// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1766",
		Title:    "Error on `memcached -l 0.0.0.0` — memcached exposed on every interface",
		Severity: SeverityError,
		Description: "`memcached -l 0.0.0.0` (or `::`, `--listen=0.0.0.0`) binds memcached's TCP " +
			"listener to every interface on the host. Memcached has no authentication and, " +
			"before `-U 0` became default, its UDP handler was the largest DDoS-" +
			"amplification vector on the internet. Bind to `127.0.0.1` or a private-" +
			"network IP only, and put memcached behind a firewall / security group scoped " +
			"to the application that consumes it.",
		Check: checkZC1766,
	})
}

func checkZC1766(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "memcached" {
		return nil
	}

	prevListen := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevListen {
			if zc1766IsUnrestrictedBind(v) {
				return zc1766Hit(cmd, "-l "+v)
			}
			prevListen = false
			continue
		}
		switch {
		case v == "-l":
			prevListen = true
		case strings.HasPrefix(v, "-l") && len(v) > 2:
			if zc1766IsUnrestrictedBind(v[2:]) {
				return zc1766Hit(cmd, v)
			}
		case strings.HasPrefix(v, "--listen="):
			if zc1766IsUnrestrictedBind(strings.TrimPrefix(v, "--listen=")) {
				return zc1766Hit(cmd, v)
			}
		}
	}
	return nil
}

func zc1766IsUnrestrictedBind(s string) bool {
	return s == "0.0.0.0" || s == "::" || s == "[::]"
}

func zc1766Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1766",
		Message: "`memcached " + what + "` exposes the unauthenticated cache to every " +
			"interface on the host. Bind to `127.0.0.1` or a private-network IP and " +
			"firewall the port.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
