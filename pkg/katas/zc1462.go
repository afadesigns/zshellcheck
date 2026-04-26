// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1462",
		Title:    "Avoid `docker run --ipc=host` — shares host IPC namespace (/dev/shm, SysV IPC)",
		Severity: SeverityWarning,
		Description: "`--ipc=host` makes the container share `/dev/shm` and the SysV IPC keyspace " +
			"with the host. Any process on the host can read/write the container's shared memory " +
			"(and vice-versa), making side-channel and data-theft attacks trivial. Use the default " +
			"private IPC namespace unless two containers explicitly need to share IPC.",
		Check: checkZC1462,
	})
}

func checkZC1462(node ast.Node) []Violation {
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

	var prevIpc bool
	for _, arg := range cmd.Arguments {
		v := arg.String()

		if v == "--ipc=host" {
			return violateZC1462(cmd)
		}
		if prevIpc {
			prevIpc = false
			if v == "host" {
				return violateZC1462(cmd)
			}
		}
		if v == "--ipc" {
			prevIpc = true
		}
	}

	return nil
}

func violateZC1462(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1462",
		Message: "`--ipc=host` shares host shared memory and SysV IPC with the container — " +
			"trivial data theft and side-channel vector. Use the default private IPC.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
