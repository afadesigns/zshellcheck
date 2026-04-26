// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1545",
		Title:    "Warn on `docker system prune -af --volumes` — drops unused volumes too",
		Severity: SeverityWarning,
		Description: "`docker system prune -af --volumes` removes stopped containers, unused " +
			"networks, dangling images — and every volume not currently attached to a running " +
			"container. On a host where `docker-compose down` is used casually (shutdown " +
			"before a laptop close, for example), the matching database volume looks " +
			"\"unused\" to prune and goes with it. Drop `--volumes` from routine cleanup, or " +
			"target specific prune scopes (`docker image prune`, `docker container prune`).",
		Check: checkZC1545,
	})
}

func checkZC1545(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" && ident.Value != "nerdctl" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	if len(args) < 2 {
		return nil
	}
	// docker system prune / volume prune
	if !((args[0] == "system" && args[1] == "prune") ||
		(args[0] == "volume" && args[1] == "prune")) {
		return nil
	}

	var hasAllVolumes bool
	for _, a := range args[2:] {
		if a == "--volumes" || a == "-a" || a == "--all" ||
			a == "-af" || a == "-fa" || a == "--all --volumes" {
			hasAllVolumes = true
		}
	}
	if !hasAllVolumes {
		return nil
	}
	return []Violation{{
		KataID: "ZC1545",
		Message: "`" + ident.Value + " " + args[0] + " prune` with `-a`/`--volumes` drops " +
			"unused volumes — stopped stacks lose their databases. Scope the prune.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
