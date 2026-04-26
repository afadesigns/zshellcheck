// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1996",
		Title:    "Warn on `unshare -U` / `-r` — unprivileged user namespace maps caller to root inside the NS",
		Severity: SeverityWarning,
		Description: "`unshare -U` opens a new user namespace and `-r` / `--map-root-user` " +
			"maps the caller's UID to `0` inside it. That's the foundation of rootless " +
			"containers (bubblewrap, podman rootless, flatpak) and is legitimate in " +
			"that context. It is also the standard opening move for a long list of " +
			"LPE chains — once you are uid `0` in a user namespace you can create " +
			"additional mount/net/cgroup namespaces, run `mount -t overlay` against " +
			"attacker-controlled dirs, and probe kernel attack surface that is " +
			"normally gated on `CAP_SYS_ADMIN`. Audit rules should flag the pattern " +
			"in production scripts; if a rootless runtime really needs it, route " +
			"through the runtime binary (`bwrap`, `podman --rootless`) so the invocation " +
			"is recognisable.",
		Check: checkZC1996,
	})
}

var zc1996ExactFlags = map[string]struct{}{
	"-U": {}, "-r": {}, "-Ur": {}, "-rU": {},
	"--user": {}, "--map-root-user": {},
}

func checkZC1996(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || CommandIdentifier(cmd) != "unshare" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1996FlagHit(v) {
			return zc1996Hit(cmd, "unshare "+v)
		}
	}
	return nil
}

func zc1996FlagHit(v string) bool {
	if _, hit := zc1996ExactFlags[v]; hit {
		return true
	}
	// Short-flag bundles like `-Urm` carry user / root mapping via U or r.
	if !strings.HasPrefix(v, "-") || strings.HasPrefix(v, "--") || len(v) <= 1 || strings.Contains(v, "=") {
		return false
	}
	return strings.ContainsAny(v[1:], "Ur")
}

func zc1996Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1996",
		Message: "`" + form + "` opens a user namespace and maps the caller to uid 0 " +
			"inside it — also the standard opening move for many kernel-LPE " +
			"chains. Route legit rootless runtimes through `bwrap`/`podman " +
			"--rootless` so the intent is clear.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
