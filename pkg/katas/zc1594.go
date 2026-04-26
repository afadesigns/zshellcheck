// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1594",
		Title:    "Warn on `docker/podman run --security-opt=systempaths=unconfined` — unhides host kernel knobs",
		Severity: SeverityWarning,
		Description: "`systempaths=unconfined` removes the container runtime's masking of " +
			"`/proc/sys`, `/proc/sysrq-trigger`, `/sys/firmware`, and related kernel surfaces. " +
			"Without the default shield a compromised process inside the container can write " +
			"`/proc/sysrq-trigger` to panic the host, or edit `/proc/sys/kernel/*` to change " +
			"kernel policy on the fly. Keep the default `systempaths=all` (masked) unless you " +
			"have a specific kernel tunable you need, then mount only that path.",
		Check: checkZC1594,
	})
}

func checkZC1594(node ast.Node) []Violation {
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

	runSeen := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !runSeen {
			if v == "run" || v == "create" {
				runSeen = true
			}
			continue
		}
		if strings.Contains(v, "systempaths=unconfined") {
			return []Violation{{
				KataID: "ZC1594",
				Message: "`" + ident.Value + " run --security-opt=systempaths=unconfined` " +
					"unhides `/proc/sys`, `/proc/sysrq-trigger`, and other kernel knobs. A " +
					"compromise in the container can then panic or re-tune the host.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
