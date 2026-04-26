// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1426",
		Title:    "Avoid `git clone http://` — unencrypted transport, use `https://` or `git://`+verify",
		Severity: SeverityWarning,
		Description: "`git clone http://...` transfers repository content unencrypted and " +
			"unauthenticated — susceptible to MITM insertion of malicious commits. Use " +
			"`https://` for authenticated hosts (GitHub, GitLab) or SSH (`git@host:path`) with " +
			"verified host keys. Plain `http://` has no integrity guarantee.",
		Check: checkZC1426,
	})
}

func checkZC1426(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	isClone := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "clone" {
			isClone = true
			continue
		}
		if isClone && strings.HasPrefix(v, "http://") {
			return []Violation{{
				KataID: "ZC1426",
				Message: "`git clone http://` is unencrypted/unauthenticated. Use `https://` " +
					"or SSH with verified host keys.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
