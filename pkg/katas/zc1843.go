// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1843",
		Title:    "Warn on `docker/podman run --cgroup-parent=/system.slice|/init.scope|/` — container escapes engine limits",
		Severity: SeverityWarning,
		Description: "`--cgroup-parent=PATH` places the container under the given cgroup parent, " +
			"which is normally `/docker` (or the engine's managed slice) and inherits the " +
			"engine-wide memory/CPU/IO caps. Pointing the flag at `/`, `/system.slice`, or " +
			"any host-managed slice puts the container side-by-side with systemd services — " +
			"the engine's defaults no longer apply, and a runaway container can starve " +
			"`sshd` or the kubelet for resources. Unless a specific orchestrator is " +
			"supplying a managed cgroup path, drop the flag and let the engine choose; if " +
			"you need custom limits, use `--memory` / `--cpus` / `--pids-limit` on the run " +
			"itself.",
		Check: checkZC1843,
	})
}

func checkZC1843(node ast.Node) []Violation {
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
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}
	sub := args[0].String()
	if sub != "run" && sub != "create" {
		return nil
	}

	for i := 1; i < len(args); i++ {
		v := args[i].String()
		var parent string
		switch {
		case strings.HasPrefix(v, "--cgroup-parent="):
			parent = strings.TrimPrefix(v, "--cgroup-parent=")
		case v == "--cgroup-parent" && i+1 < len(args):
			parent = args[i+1].String()
		default:
			continue
		}
		if zc1843IsHostSlice(parent) {
			return []Violation{{
				KataID: "ZC1843",
				Message: "`" + ident.Value + " " + sub + " --cgroup-parent=" + parent +
					"` puts the container under a host-managed slice — the engine's " +
					"memory/CPU caps no longer apply. Drop the flag or pass " +
					"`--memory`/`--cpus`/`--pids-limit` directly.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1843IsHostSlice(v string) bool {
	if v == "" {
		return false
	}
	trimmed := strings.Trim(v, "\"'")
	switch trimmed {
	case "/", "/system.slice", "/user.slice", "/init.scope", "/machine.slice":
		return true
	}
	// Anything under /system.slice or /init.scope qualifies too.
	for _, prefix := range []string{"/system.slice/", "/init.scope/", "/machine.slice/"} {
		if strings.HasPrefix(trimmed, prefix) {
			return true
		}
	}
	return false
}
