// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1658",
		Title:    "Warn on `curl -OJ` / `-J -O` — server-controlled output filename",
		Severity: SeverityWarning,
		Description: "`curl -J` (`--remote-header-name`) combined with `-O` (`--remote-name`) " +
			"saves the response using the filename the server puts in the `Content-Disposition` " +
			"header. The server — or anything on the path that can set headers, including a " +
			"compromised CDN or an HTTP-serving reverse proxy — chooses the destination name. " +
			"Paths like `../../etc/cron.d/evil` are rejected by curl's sanitizer, but benign-" +
			"looking names still overwrite files in the current directory. Use `-o NAME` with " +
			"a filename you control, and validate the payload before you act on it.",
		Check: checkZC1658,
	})
}

func checkZC1658(node ast.Node) []Violation {
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

	hasJ := false
	hasO := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--remote-header-name" {
			hasJ = true
			continue
		}
		if v == "--remote-name" {
			hasO = true
			continue
		}
		if !strings.HasPrefix(v, "-") || strings.HasPrefix(v, "--") {
			continue
		}
		body := strings.TrimPrefix(v, "-")
		if strings.ContainsRune(body, 'J') {
			hasJ = true
		}
		if strings.ContainsRune(body, 'O') {
			hasO = true
		}
	}

	if !hasJ || !hasO {
		return nil
	}

	return []Violation{{
		KataID: "ZC1658",
		Message: "`curl -OJ` saves the response under the name the server picks in " +
			"`Content-Disposition` — path traversal is blocked but arbitrary same-dir " +
			"overwrites are not. Pass `-o NAME` with a filename you control.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
