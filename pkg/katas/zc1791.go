// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1791DaemonSockets = []string{
	"/var/run/docker.sock",
	"/run/docker.sock",
	"/var/run/podman/podman.sock",
	"/run/podman/podman.sock",
	"/run/containerd/containerd.sock",
	"/run/crio/crio.sock",
	"/var/run/docker/containerd/containerd.sock",
}

var zc1791UnixSocketFlags = map[string]bool{"--unix-socket": true}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1791",
		Title:    "Error on `curl --unix-socket /var/run/docker.sock` — direct container-daemon API access",
		Severity: SeverityError,
		Description: "A curl request to `docker.sock` / `containerd.sock` / `crio.sock` speaks " +
			"the container-daemon HTTP API with no authentication beyond the socket's " +
			"filesystem permissions. Anyone who can invoke curl as that uid can `POST " +
			"/containers/create` with `HostConfig.Privileged=true` and a bind mount of `/` " +
			"and land a root shell on the host — the primitive every \"docker socket " +
			"escape\" write-up leans on. Use the real CLI (`docker`, `podman`, `nerdctl`) " +
			"which enforces its own policy, or access the daemon over a TLS-protected TCP " +
			"endpoint with mutual auth.",
		Check: checkZC1791,
	})
}

func checkZC1791(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "curl" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		var path string
		switch {
		case v == "--unix-socket":
			if i+1 >= len(cmd.Arguments) {
				return nil
			}
			path = cmd.Arguments[i+1].String()
		case strings.HasPrefix(v, "--unix-socket="):
			path = strings.TrimPrefix(v, "--unix-socket=")
		default:
			continue
		}
		path = strings.Trim(path, "\"'")
		if hit := zc1791MatchSocket(cmd, path); hit != nil {
			return hit
		}
	}
	return nil
}

func zc1791MatchSocket(cmd *ast.SimpleCommand, path string) []Violation {
	for _, sock := range zc1791DaemonSockets {
		if path == sock {
			line, col := FlagArgPosition(cmd, zc1791UnixSocketFlags)
			return []Violation{{
				KataID: "ZC1791",
				Message: "`curl --unix-socket " + path + "` speaks the container-daemon " +
					"API — a `POST /containers/create` with `Privileged=true` is a " +
					"host-root primitive. Use the CLI (`docker`/`podman`) instead.",
				Line:   line,
				Column: col,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
