// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1651",
		Title:    "Warn on `docker/podman run -p 0.0.0.0:PORT:PORT` — explicit all-interfaces publish",
		Severity: SeverityWarning,
		Description: "A port spec of `0.0.0.0:HOST:CONT`, `[::]:HOST:CONT`, or `*:HOST:CONT` " +
			"publishes the container port to every interface the host has. On a multi-" +
			"tenant LAN or a cloud host with a public IP the service is immediately reachable " +
			"from anywhere. If the service needs only local reverse-proxy access, bind to " +
			"`127.0.0.1:HOST:CONT` and let nginx / caddy handle external exposure.",
		Check: checkZC1651,
	})
}

func checkZC1651(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "run" && sub != "create" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v != "-p" && v != "--publish" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		spec := strings.Trim(cmd.Arguments[i+1].String(), "\"'")
		if strings.HasPrefix(spec, "0.0.0.0:") ||
			strings.HasPrefix(spec, "[::]:") ||
			strings.HasPrefix(spec, "*:") {
			return []Violation{{
				KataID: "ZC1651",
				Message: "`" + ident.Value + " " + sub + " -p " + spec + "` publishes to " +
					"every interface. Bind to `127.0.0.1:HOST:CONT` and put nginx / caddy " +
					"in front for external access.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
