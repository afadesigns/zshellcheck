// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1884",
		Title:    "Error on `curl/wget https://...?apikey=...` — credential in URL query string",
		Severity: SeverityError,
		Description: "Anything passed as an HTTP query parameter is logged by every intermediary: " +
			"the server's access log, the transparent proxy, the CDN request-id trail, " +
			"browser referrer headers, and any client-side observability tooling. A URL " +
			"like `https://api.example/widgets?apikey=SECRET&token=xyz` therefore " +
			"tattoos the credential into logs that live forever and are often shared " +
			"with downstream teams. Move the secret into an HTTP header " +
			"(`curl -H \"Authorization: Bearer $TOKEN\"`), a POST body with " +
			"`--data-urlencode` + TLS, or an `-u user:` basic-auth combo — never the " +
			"query string.",
		Check: checkZC1884,
	})
}

var zc1884SecretKeys = []string{
	"apikey=",
	"api_key=",
	"api-key=",
	"token=",
	"access_token=",
	"id_token=",
	"auth_token=",
	"access-token=",
	"password=",
	"passwd=",
	"secret=",
	"client_secret=",
	"sig=",
	"signature=",
}

func checkZC1884(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "curl" && ident.Value != "wget" && ident.Value != "http" && ident.Value != "httpie" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if match := zc1884FirstSecretKey(v); match != "" {
			return []Violation{{
				KataID: "ZC1884",
				Message: "`" + ident.Value + " " + v + "` carries `" + match +
					"...` in the URL query — logged by every proxy, CDN, and " +
					"server access log along the path. Move credentials to " +
					"`-H \"Authorization: Bearer \"$TOKEN\"` or a POST body.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1884FirstSecretKey(v string) string {
	lower := strings.ToLower(strings.Trim(v, "\"'"))
	// Only flag when this looks like a URL; drop anything without `://` or `?`.
	if !strings.Contains(lower, "://") || !strings.Contains(lower, "?") {
		return ""
	}
	// Scan after the `?` boundary.
	idx := strings.Index(lower, "?")
	if idx < 0 {
		return ""
	}
	query := lower[idx:]
	for _, key := range zc1884SecretKeys {
		if strings.Contains(query, key) {
			return strings.TrimSuffix(key, "=")
		}
	}
	return ""
}
