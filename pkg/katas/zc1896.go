// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1896",
		Title:    "Error on `docker/podman run -v /proc:…|/sys:…` — bind-mounts host kernel interfaces into container",
		Severity: SeverityError,
		Description: "`docker run -v /proc:/host/proc` (or `-v /sys:…`) bind-mounts the host's " +
			"procfs / sysfs hierarchy into the container's mount namespace. From inside, " +
			"the container can read every host process's `environ` (secrets passed via " +
			"env), every `cmdline`, every `/proc/1/ns/` to open namespace fds for a " +
			"breakout, and `/sys/fs/cgroup` to modify resource limits that affect host " +
			"services. `:ro` does not help — `/proc/<pid>/ns/...` handles remain usable. " +
			"If the container genuinely needs process / kernel visibility, grant the " +
			"narrowest capability instead (`--cap-add=SYS_PTRACE`) or run the monitoring " +
			"agent on the host rather than inside an untrusted image.",
		Check: checkZC1896,
	})
}

func checkZC1896(node ast.Node) []Violation {
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
	if len(args) == 0 || (args[0].String() != "run" && args[0].String() != "create") {
		return nil
	}
	for i := 1; i < len(args); i++ {
		v := args[i].String()
		var mount string
		switch {
		case (v == "-v" || v == "--volume" || v == "--mount") && i+1 < len(args):
			mount = args[i+1].String()
		case strings.HasPrefix(v, "--volume=") || strings.HasPrefix(v, "--mount=") || strings.HasPrefix(v, "-v="):
			mount = v[strings.Index(v, "=")+1:]
		default:
			continue
		}
		if src := zc1896HostKernelSource(mount); src != "" {
			return []Violation{{
				KataID: "ZC1896",
				Message: "`" + ident.Value + " ... -v " + mount + "` bind-mounts host " +
					src + " into the container — every process's `environ`/`cmdline` " +
					"and `/proc/1/ns/` breakout handles become readable. Use " +
					"`--cap-add=SYS_PTRACE` or host-side monitoring instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1896HostKernelSource(v string) string {
	trimmed := strings.Trim(v, "\"'")
	// Accept `source:target[:opts]` bind form and `source=/path,…` mount form.
	source := trimmed
	if idx := strings.Index(trimmed, ":"); idx > 0 {
		source = trimmed[:idx]
	}
	// `--mount type=bind,source=/proc,…`
	if strings.Contains(trimmed, "source=") {
		for _, entry := range strings.Split(trimmed, ",") {
			if strings.HasPrefix(entry, "source=") {
				source = strings.TrimPrefix(entry, "source=")
			} else if strings.HasPrefix(entry, "src=") {
				source = strings.TrimPrefix(entry, "src=")
			}
		}
	}
	switch source {
	case "/proc", "/sys":
		return source
	}
	if strings.HasPrefix(source, "/proc/") || strings.HasPrefix(source, "/sys/") {
		return source
	}
	return ""
}
