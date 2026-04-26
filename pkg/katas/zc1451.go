// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1451",
		Title:    "Avoid `pip install` without `--user` or virtualenv",
		Severity: SeverityWarning,
		Description: "`pip install pkg` (no `--user`, no active venv) targets the system Python, " +
			"potentially breaking system tools or requiring sudo. On modern Linux this now fails " +
			"with PEP 668 `externally-managed-environment`. Always use a virtualenv (`python -m " +
			"venv`, `uv`, `poetry`) or `--user` for scoped installs.",
		Check: checkZC1451,
	})
}

func checkZC1451(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pip" && ident.Value != "pip3" {
		return nil
	}

	hasInstall := false
	hasUser := false
	hasBreakSystem := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" {
			hasInstall = true
		}
		if v == "--user" {
			hasUser = true
		}
		if v == "--break-system-packages" {
			hasBreakSystem = true
		}
	}
	if hasInstall && !hasUser && !hasBreakSystem {
		return []Violation{{
			KataID: "ZC1451",
			Message: "`pip install` without `--user` or an active venv targets system Python. " +
				"Use `python -m venv` / `uv` / `--user` for scoped installs.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
