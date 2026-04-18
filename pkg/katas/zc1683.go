package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1683",
		Title:    "Error on `npm/yarn/pnpm config set registry http://...` — plaintext package index",
		Severity: SeverityError,
		Description: "Pointing a JavaScript package manager at an `http://` registry disables " +
			"TLS during fetch. Any host on the path (corporate proxy, hotel Wi-Fi, " +
			"compromised CDN) can rewrite tarballs mid-flight; lockfile hashes catch the " +
			"rewrite only if the user locks every dependency before the swap. Even on " +
			"internal networks, pin to `https://` — reach for your own CA via " +
			"`NODE_EXTRA_CA_CERTS` or `registry.cafile` rather than falling back to HTTP.",
		Check: checkZC1683,
	})
}

func checkZC1683(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "npm" && ident.Value != "yarn" && ident.Value != "pnpm" {
		return nil
	}

	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "config" || cmd.Arguments[1].String() != "set" {
		return nil
	}
	if cmd.Arguments[2].String() != "registry" {
		return nil
	}
	url := cmd.Arguments[3].String()
	if !strings.HasPrefix(url, "http://") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1683",
		Message: "`" + ident.Value + " config set registry " + url + "` uses plaintext " +
			"HTTP — any proxy / CDN can rewrite tarballs. Use `https://` and a custom " +
			"CA via `NODE_EXTRA_CA_CERTS` if needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
