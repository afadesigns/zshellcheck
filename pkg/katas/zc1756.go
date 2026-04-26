// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1756DockerSocketPaths = map[string]bool{
	"/var/run/docker.sock":              true,
	"/run/docker.sock":                  true,
	"/run/containerd/containerd.sock":   true,
	"/var/run/crio/crio.sock":           true,
	"/run/crio/crio.sock":               true,
	"/var/run/podman/podman.sock":       true,
	"/run/podman/podman.sock":           true,
	"/run/user/1000/podman/podman.sock": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1756",
		Title:    "Error on `chmod NNN /run/docker.sock` — world access is root-equivalent privesc",
		Severity: SeverityError,
		Description: "Container-runtime sockets (`/var/run/docker.sock`, `/run/containerd/" +
			"containerd.sock`, `/run/crio/crio.sock`, `/run/podman/podman.sock`) accept " +
			"commands that run on the host with root privilege — starting privileged " +
			"containers, mounting the host filesystem, reading every file on disk. " +
			"Making the socket world-readable or world-writable (`chmod 644/660/666/777`) " +
			"hands every local user that root-escalation primitive. Keep the socket " +
			"`0660 root:docker` (or the equivalent runtime group) and add only trusted " +
			"accounts to that group.",
		Check: checkZC1756,
	})
}

func checkZC1756(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	var mode, target string
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1756DockerSocketPaths[v] {
			target = v
			continue
		}
		if mode == "" && zc1756WorldAccess(v) {
			mode = v
		}
	}
	if mode == "" || target == "" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1756",
		Message: "`chmod " + mode + " " + target + "` grants every local user access to a " +
			"root-equivalent container-runtime socket. Keep `0660` owned by the " +
			"runtime group (`root:docker` etc.) and restrict membership.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

// zc1756WorldAccess: true when MODE grants world-read or world-write on the
// target. Handles both octal literal and the parser-normalised decimal form.
func zc1756WorldAccess(mode string) bool {
	if mode == "" {
		return false
	}
	if zc1671WorldWritable(mode) {
		return true
	}
	// zc1671WorldWritable catches world-write; also catch world-read (bit 4).
	// Parse octal, fallback to decimal.
	for _, c := range mode {
		if c < '0' || c > '7' {
			// Decimal fallback (parser-normalised leading-zero octal).
			return false
		}
	}
	// Simple octal parse: last char is world triad.
	last := mode[len(mode)-1]
	return last == '4' || last == '5' || last == '6' || last == '7'
}
