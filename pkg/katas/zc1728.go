// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1728",
		Title:    "Error on `pip install --index-url http://...` — plaintext index allows MITM",
		Severity: SeverityError,
		Description: "`pip install --index-url http://...`, `--extra-index-url http://...`, " +
			"and `-i http://...` tell pip to fetch packages over plaintext HTTP. Any " +
			"network-position attacker (open Wi-Fi, hostile transit, MITM proxy) can " +
			"replace package metadata or wheel contents in flight — direct code execution " +
			"on the install host. Switch to `https://`, or on internal networks terminate " +
			"TLS at the mirror and only configure the `https://` URL.",
		Check: checkZC1728,
	})
}

func checkZC1728(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "pip", "pip3", "pip2":
	default:
		return nil
	}

	prevURL := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevURL {
			if zc1728PlainHTTP(v) {
				return zc1728Hit(cmd, v)
			}
			prevURL = false
			continue
		}
		switch {
		case v == "--index-url" || v == "--extra-index-url" || v == "-i":
			prevURL = true
		case strings.HasPrefix(v, "--index-url="), strings.HasPrefix(v, "--extra-index-url="):
			eq := strings.IndexByte(v, '=')
			if zc1728PlainHTTP(v[eq+1:]) {
				return zc1728Hit(cmd, v)
			}
		}
	}
	return nil
}

func zc1728PlainHTTP(url string) bool {
	return strings.HasPrefix(url, "http://")
}

func zc1728Hit(cmd *ast.SimpleCommand, url string) []Violation {
	return []Violation{{
		KataID: "ZC1728",
		Message: "`pip install --index-url " + url + "` fetches packages over plaintext " +
			"HTTP — any MITM swaps the wheel for code execution on the host. Use " +
			"`https://`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
