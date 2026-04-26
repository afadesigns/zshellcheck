// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1463",
		Title:    "Avoid `docker run --userns=host` — disables user-namespace remapping",
		Severity: SeverityWarning,
		Description: "`--userns=host` turns off the user-namespace remap, meaning UID 0 in the " +
			"container maps to UID 0 on the host. Combined with any of the `--cap-add`, " +
			"`--privileged`, or bind-mount footguns, this becomes a direct host-root escalation. " +
			"Leave the default (container-side remap) enabled.",
		Check: checkZC1463,
	})
}

func checkZC1463(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	var prevNs bool
	for _, arg := range cmd.Arguments {
		v := arg.String()

		if v == "--userns=host" {
			return violateZC1463(cmd)
		}
		if prevNs {
			prevNs = false
			if v == "host" {
				return violateZC1463(cmd)
			}
		}
		if v == "--userns" {
			prevNs = true
		}
	}

	return nil
}

func violateZC1463(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1463",
		Message: "`--userns=host` disables user-namespace remap — UID 0 in the container == " +
			"UID 0 on the host. Leave the default remap on.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
