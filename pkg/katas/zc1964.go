// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1964",
		Title:    "Warn on `uvx pkg` / `uv tool run pkg` / `pipx run pkg` without a version pin — runs latest PyPI release",
		Severity: SeverityWarning,
		Description: "`uvx PKG`, `uv tool run PKG`, and `pipx run PKG` each resolve the package " +
			"against PyPI and execute its entry point. Without a version constraint " +
			"(`pkg==1.2.3` or `pkg@1.2.3` for uv), every run takes whatever the registry " +
			"currently serves — a typosquatted lookalike, a compromised maintainer " +
			"release, or a sudden major-version bump lands untested code in the " +
			"pipeline. Pin the version at the call site or use `uv tool install pkg==X.Y.Z` + " +
			"`uv tool run pkg` so the lockfile is the source of truth.",
		Check: checkZC1964,
	})
}

func checkZC1964(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var form string
	var pkgs []ast.Expression
	switch ident.Value {
	case "uvx":
		form = "uvx"
		pkgs = cmd.Arguments
	case "uv":
		if len(cmd.Arguments) < 3 {
			return nil
		}
		if cmd.Arguments[0].String() != "tool" || cmd.Arguments[1].String() != "run" {
			return nil
		}
		form = "uv tool run"
		pkgs = cmd.Arguments[2:]
	case "pipx":
		if len(cmd.Arguments) < 2 {
			return nil
		}
		if cmd.Arguments[0].String() != "run" {
			return nil
		}
		form = "pipx run"
		pkgs = cmd.Arguments[1:]
	default:
		return nil
	}

	for _, arg := range pkgs {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			continue
		}
		if zc1964IsPinned(v) {
			return nil
		}
		if strings.HasPrefix(v, "$") {
			return nil
		}
		return []Violation{{
			KataID: "ZC1964",
			Message: "`" + form + " " + v + "` resolves to the PyPI `latest` release — " +
				"a squatted name or compromised maintainer lands untested code. Pin " +
				"`pkg==X.Y.Z` (or `pkg@X.Y.Z` for uv).",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

func zc1964IsPinned(v string) bool {
	return strings.Contains(v, "==") || strings.Contains(v, "@") ||
		strings.Contains(v, ">=") || strings.Contains(v, "~=")
}
