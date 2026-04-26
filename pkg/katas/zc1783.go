// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1783",
		Title:    "Error on `podman system reset` / `nerdctl system prune -af --volumes` — wipes every container artifact",
		Severity: SeverityError,
		Description: "`podman system reset` removes every podman container, image, volume, " +
			"network, pod, secret, and storage driver scratch area — a full factory reset " +
			"of the local engine. `nerdctl system prune -af --volumes` achieves the same for " +
			"containerd. On a developer workstation this wipes cached images for unrelated " +
			"projects; on a CI runner or build host it invalidates every warm artifact the " +
			"job relies on; on a prod host it drops the volumes the workload stores data in. " +
			"Use narrower commands (`podman rmi`, `podman volume rm`, scoped `podman prune`) " +
			"that only touch the resource you intend to remove, and never pair the reset with " +
			"`--force`.",
		Check: checkZC1783,
	})
}

func checkZC1783(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "podman":
		if len(cmd.Arguments) >= 2 &&
			cmd.Arguments[0].String() == "system" &&
			cmd.Arguments[1].String() == "reset" {
			return zc1783Hit(cmd, "podman system reset")
		}
	case "nerdctl":
		if len(cmd.Arguments) >= 2 &&
			cmd.Arguments[0].String() == "system" &&
			cmd.Arguments[1].String() == "prune" {
			hasAll := false
			hasVolumes := false
			for _, arg := range cmd.Arguments[2:] {
				v := arg.String()
				if v == "-af" || v == "-fa" || v == "-a" || v == "--all" {
					hasAll = true
				}
				if v == "--volumes" {
					hasVolumes = true
				}
			}
			if hasAll && hasVolumes {
				return zc1783Hit(cmd, "nerdctl system prune -a --volumes")
			}
		}
	}
	return nil
}

func zc1783Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1783",
		Message: "`" + what + "` wipes every container artifact on the host — images, " +
			"volumes, networks, pods. Use narrower removals (`rmi`, `volume rm`, scoped " +
			"`prune`) against the specific resource you intend to delete.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
