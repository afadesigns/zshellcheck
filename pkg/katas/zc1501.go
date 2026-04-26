// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1501",
		Title:    "Style: `docker-compose` (hyphen) — use `docker compose` (space, built-in plugin)",
		Severity: SeverityStyle,
		Description: "`docker-compose` is the Python Compose V1 binary. Docker stopped shipping " +
			"it with Docker Desktop in 2023 and Compose V2 is now the first-class `docker " +
			"compose` subcommand. Scripts that invoke `docker-compose` silently degrade on " +
			"fresh installs and miss V2-only options (`--profile`, `--wait`, richer env " +
			"interpolation). Call `docker compose` (space) or pin the V2 binary explicitly.",
		Check: checkZC1501,
		Fix:   fixZC1501,
	})
}

// fixZC1501 rewrites the hyphenated `docker-compose` command name into
// the space-separated `docker compose` subcommand form. Arguments stay
// untouched — the V2 plugin accepts the same shape.
func fixZC1501(node ast.Node, v Violation, _ []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker-compose" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("docker-compose"),
		Replace: "docker compose",
	}}
}

func checkZC1501(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker-compose" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1501",
		Message: "`docker-compose` is the deprecated Python V1 binary. Use `docker compose` " +
			"(space-separated subcommand) for the bundled V2 plugin.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
