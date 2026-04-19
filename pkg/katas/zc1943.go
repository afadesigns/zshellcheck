package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1943",
		Title:    "Warn on `systemd-nspawn -b` / `--boot` — runs a full init inside a possibly untrusted rootfs",
		Severity: SeverityWarning,
		Description: "`systemd-nspawn -b -D $ROOT` (and `--boot -D $ROOT`) launches the rootfs's " +
			"`/sbin/init` inside a minimally-isolated namespace — by default the container " +
			"inherits `CAP_AUDIT_CONTROL`, `CAP_NET_ADMIN`, and read-write access to the " +
			"host's `/dev` nodes that match the container's cgroup. If `$ROOT` is an " +
			"operator-supplied tarball, any init script it ships runs first and can probe the " +
			"host. Use `-U` for user-namespace isolation, drop capabilities with " +
			"`--capability=`, pair with `--private-network`, and prefer `machinectl start` on " +
			"a reviewed image instead of ad-hoc boots.",
		Check: checkZC1943,
	})
}

func checkZC1943(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `systemd-nspawn --boot` mangles the command name to
	// `boot` (after the `systemd-nspawn--` subtract).
	if ident.Value == "boot" {
		return zc1943Hit(cmd, "systemd-nspawn --boot")
	}
	if ident.Value != "systemd-nspawn" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-b" || v == "--boot" {
			return zc1943Hit(cmd, "systemd-nspawn "+v)
		}
	}
	return nil
}

func zc1943Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1943",
		Message: "`" + form + "` runs the rootfs's `/sbin/init` with minimal isolation — " +
			"init scripts execute first and can probe the host. Use `-U`, drop caps with " +
			"`--capability=`, pair with `--private-network`, prefer `machinectl start`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
