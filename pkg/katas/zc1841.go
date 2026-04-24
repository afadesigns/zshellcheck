package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1841ProxyFlags = map[string]bool{"--proxy-insecure": true}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1841",
		Title:    "Error on `curl --proxy-insecure` — TLS verification disabled on the proxy hop",
		Severity: SeverityError,
		Description: "`curl --proxy-insecure` (alias of `-k` but scoped to the proxy leg, " +
			"introduced alongside `--proxy-cacert` in curl 7.52) tells curl to accept any " +
			"certificate presented by the HTTPS proxy that sits between the script and the " +
			"origin server. The origin TLS handshake is still validated, which makes the " +
			"issue easy to miss in review, but any box that can intercept traffic to the " +
			"proxy — a captive portal, a rogue WPAD auto-discovery, an attacker on the same " +
			"VLAN — can present its own cert and read or rewrite the tunnel contents, " +
			"including any `Authorization:` header attached to the request. Install the " +
			"proxy's CA bundle and point `--proxy-cacert` / `CURL_CA_BUNDLE` at it instead.",
		Check: checkZC1841,
	})
}

func checkZC1841(node ast.Node) []Violation {
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
	for _, arg := range cmd.Arguments {
		if arg.String() == "--proxy-insecure" {
			return zc1841Hit(cmd)
		}
	}
	return nil
}

func zc1841Hit(cmd *ast.SimpleCommand) []Violation {
	line, col := FlagArgPosition(cmd, zc1841ProxyFlags)
	return []Violation{{
		KataID: "ZC1841",
		Message: "`curl --proxy-insecure` skips TLS verification on the proxy hop — " +
			"an on-path attacker can present any cert and decrypt the tunnel " +
			"(including `Authorization:` headers). Install the proxy CA and use " +
			"`--proxy-cacert PATH`.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}
