// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1630",
		Title:    "Warn on `php -S 0.0.0.0:PORT` — PHP dev server exposes CWD to all interfaces",
		Severity: SeverityWarning,
		Description: "`php -S 0.0.0.0:PORT` starts PHP's built-in dev server listening on every " +
			"interface the host has. It serves files from the working directory (or the " +
			"docroot named after the bind) with no auth, no TLS, and minimal access logging. " +
			"The PHP docs explicitly say not to use it in production. Bind to `127.0.0.1:PORT` " +
			"for local testing and put nginx / caddy in front for anything externally exposed.",
		Check: checkZC1630,
	})
}

func checkZC1630(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "php" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-S" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			return nil
		}
		bind := cmd.Arguments[i+1].String()
		if strings.HasPrefix(bind, "0.0.0.0:") ||
			strings.HasPrefix(bind, "*:") ||
			strings.HasPrefix(bind, "[::]:") {
			return []Violation{{
				KataID: "ZC1630",
				Message: "`php -S " + bind + "` binds the dev server to every interface — " +
					"unauthenticated access to the working directory. Use `127.0.0.1:PORT` " +
					"locally, nginx / caddy for external exposure.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
