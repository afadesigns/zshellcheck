// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1456",
		Title:    "Avoid `docker run -v /:...` — bind-mounts host root into container",
		Severity: SeverityError,
		Description: "Mounting `/` (host root) into a container gives the container read/write " +
			"access to the entire host filesystem — a trivial container escape. Mount only the " +
			"specific host paths the container needs, using `:ro` for read-only where possible.",
		Check: checkZC1456,
	})
}

func checkZC1456(node ast.Node) []Violation {
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

	var prevV bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevV {
			prevV = false
			// v should be host:container or host:container:opts
			if strings.HasPrefix(v, "/:") || v == "/" {
				return []Violation{{
					KataID: "ZC1456",
					Message: "`-v /:...` mounts the host root into the container — trivial " +
						"container escape. Scope to specific paths.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
		if v == "-v" || v == "--volume" {
			prevV = true
		}
	}

	return nil
}
