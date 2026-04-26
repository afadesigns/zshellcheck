// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1498",
		Title:    "Warn on `mount -o remount,rw /` — makes read-only root filesystem writable",
		Severity: SeverityWarning,
		Description: "Remounting the root filesystem read-write is either an intentional config " +
			"change that belongs in `/etc/fstab` (in which case this script is the wrong place) " +
			"or a post-compromise step for persisting changes on an immutable / verity-backed " +
			"root. On distros that ship with RO root (Fedora Silverblue, Chrome OS, appliance " +
			"images) this also breaks rollback guarantees. Use `systemd-sysext` or " +
			"`ostree admin deploy` for legitimate modifications.",
		Check: checkZC1498,
	})
}

var zc1498Roots = map[string]struct{}{"/": {}, "/root": {}, "/boot": {}}

func checkZC1498(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "mount" {
		return nil
	}
	args := zc1464StringArgs(cmd)
	hasRemount, hasRW := zc1498RemountFlags(args)
	target := zc1498SystemTarget(args)
	if !hasRemount || !hasRW || target == "" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1498",
		Message: "`mount -o remount,rw " + target + "` makes a read-only system path " +
			"writable — use ostree / systemd-sysext or fix /etc/fstab.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func zc1498RemountFlags(args []string) (hasRemount, hasRW bool) {
	for i, a := range args {
		if a != "-o" || i+1 >= len(args) {
			continue
		}
		for _, o := range strings.Split(args[i+1], ",") {
			switch o {
			case "remount":
				hasRemount = true
			case "rw":
				hasRW = true
			}
		}
	}
	return
}

func zc1498SystemTarget(args []string) string {
	for _, a := range args {
		if _, hit := zc1498Roots[a]; hit && !strings.HasPrefix(a, "-") {
			return a
		}
	}
	return ""
}
