// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1610",
		Title:    "Warn on `curl -o /etc/...` / `wget -O /etc/...` — direct download to a system path",
		Severity: SeverityWarning,
		Description: "Writing the body of an HTTP response straight into `/etc/`, `/usr/`, " +
			"`/bin/`, `/sbin/`, or `/lib/` skips every integrity check the system usually " +
			"applies. If the URL is compromised or MITM'd, the attacker's content replaces a " +
			"system config or binary the next command over. Download to a temp file, verify " +
			"signature / checksum, and `install -m 0644` the final file into place. Package " +
			"managers exist for a reason — prefer them for system files.",
		Check: checkZC1610,
	})
}

func checkZC1610(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "curl" && ident.Value != "wget" {
		return nil
	}

	systemPrefixes := []string{"/etc/", "/usr/", "/bin/", "/sbin/", "/lib/", "/lib64/", "/opt/"}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		isOutFlag := v == "-o" || v == "-O" || v == "--output" || v == "--output-document"
		if isOutFlag && i+1 < len(cmd.Arguments) {
			next := cmd.Arguments[i+1].String()
			for _, p := range systemPrefixes {
				if strings.HasPrefix(next, p) {
					return []Violation{{
						KataID: "ZC1610",
						Message: "`" + ident.Value + " " + v + " " + next + "` writes " +
							"an HTTP response straight into a system path — a compromised " +
							"URL replaces the target. Download to a temp file, verify, " +
							"then `install` into place.",
						Line:   cmd.Token.Line,
						Column: cmd.Token.Column,
						Level:  SeverityWarning,
					}}
				}
			}
		}
		// Handle --output=/path and --output-document=/path joined forms.
		for _, prefix := range []string{"--output=", "--output-document="} {
			if strings.HasPrefix(v, prefix) {
				path := strings.TrimPrefix(v, prefix)
				for _, p := range systemPrefixes {
					if strings.HasPrefix(path, p) {
						return []Violation{{
							KataID: "ZC1610",
							Message: "`" + ident.Value + " " + v + "` writes an HTTP " +
								"response straight into a system path — a compromised " +
								"URL replaces the target. Download to a temp file, " +
								"verify, then `install` into place.",
							Line:   cmd.Token.Line,
							Column: cmd.Token.Column,
							Level:  SeverityWarning,
						}}
					}
				}
			}
		}
	}
	return nil
}
