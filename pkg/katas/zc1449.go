// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1449",
		Title:    "`dnf`/`yum` install without `-y` hangs in non-interactive scripts",
		Severity: SeverityWarning,
		Description: "In CI/Dockerfiles, `dnf install pkg` or `yum install pkg` prompts for " +
			"confirmation and stalls. Always pass `-y` (or `--assumeyes`) for unattended runs. " +
			"Also consider `--nodocs` and `--setopt=install_weak_deps=False` for slim images.",
		Check: checkZC1449,
	})
}

func checkZC1449(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dnf" && ident.Value != "yum" && ident.Value != "microdnf" {
		return nil
	}

	hasInstall := false
	hasYes := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" || v == "upgrade" || v == "update" || v == "remove" {
			hasInstall = true
		}
		if v == "-y" || v == "--assumeyes" {
			hasYes = true
		}
	}
	if hasInstall && !hasYes {
		return []Violation{{
			KataID: "ZC1449",
			Message: "`" + ident.Value + "` without `-y` hangs on confirmation. Add `-y` for " +
				"unattended runs.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
