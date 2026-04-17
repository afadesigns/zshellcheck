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

func checkZC1498(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "mount" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	var hasRemount, hasRW bool
	var target string
	for i, a := range args {
		if a == "-o" && i+1 < len(args) {
			opts := strings.Split(args[i+1], ",")
			for _, o := range opts {
				if o == "remount" {
					hasRemount = true
				}
				if o == "rw" {
					hasRW = true
				}
			}
		}
		if (a == "/" || a == "/root" || a == "/boot") && !strings.HasPrefix(a, "-") {
			target = a
		}
	}
	if hasRemount && hasRW && target != "" {
		return []Violation{{
			KataID: "ZC1498",
			Message: "`mount -o remount,rw " + target + "` makes a read-only system path " +
				"writable — use ostree / systemd-sysext or fix /etc/fstab.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
