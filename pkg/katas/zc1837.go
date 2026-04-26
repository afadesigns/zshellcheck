// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strconv"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1837",
		Title:    "Error on `chmod` granting non-owner access to `/dev/kvm` / `/dev/mem` / `/dev/kmem` / `/dev/port`",
		Severity: SeverityError,
		Description: "Distros ship `/dev/mem`, `/dev/kmem`, `/dev/port`, and `/dev/kvm` with tight " +
			"owner-only or group-only permissions managed by udev rules — these nodes hand " +
			"any process that can read or write them the keys to the kingdom (physical " +
			"memory, kernel memory, raw I/O ports, full hypervisor API). Flipping the mode " +
			"from a script (`chmod 666 /dev/kvm`, `chmod a+rw /dev/mem`) is a classic local " +
			"privilege-escalation vector dressed up as a convenience fix for a permission " +
			"error. Fix the actual problem: add the user to the `kvm` group, ship a proper " +
			"udev rule (`/etc/udev/rules.d/*.rules`), or grant the specific capability the " +
			"tool needs instead of blanket-chmod-ing the device.",
		Check: checkZC1837,
	})
}

var zc1837Devices = map[string]struct{}{
	"/dev/kvm":  {},
	"/dev/mem":  {},
	"/dev/kmem": {},
	"/dev/port": {},
}

func checkZC1837(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}

	var mode, target string
	for _, arg := range args {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			continue
		}
		if mode == "" {
			mode = v
			continue
		}
		target = v
		break
	}
	if target == "" {
		return nil
	}
	if _, hit := zc1837Devices[target]; !hit {
		return nil
	}
	if !zc1837GrantsNonOwner(mode) {
		return nil
	}
	return []Violation{{
		KataID: "ZC1837",
		Message: "`chmod " + mode + " " + target + "` grants non-owner access to " +
			"a privileged kernel device — classic local-privesc vector. Use " +
			"group membership or a udev rule instead of blanket chmod.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func zc1837GrantsNonOwner(mode string) bool {
	if mode == "" {
		return false
	}
	// Symbolic: look for tokens that grant to group/other/all.
	lower := strings.ToLower(mode)
	if strings.HasPrefix(lower, "+") {
		return true
	}
	for _, frag := range []string{"o+", "o=", "a+", "a=", "ugo+", "ugo="} {
		if strings.Contains(lower, frag) {
			return true
		}
	}
	// Numeric: chmod reads the mode as octal. Parser normalises leading-zero
	// octals to decimal (e.g. "0666" -> "438"), so branch on which one we got.
	for _, r := range mode {
		if r < '0' || r > '9' {
			return false
		}
	}
	var n int64
	if strings.ContainsAny(mode, "89") {
		n, _ = strconv.ParseInt(mode, 10, 32)
	} else {
		n, _ = strconv.ParseInt(mode, 8, 32)
	}
	// Flag only if any "other" (world) bit is set — these devices are managed
	// with group-only access (e.g. /dev/kvm = 660 kvm:root); tightening to
	// 660/600 is fine, opening to the world is the privesc case.
	return (n & 0o007) != 0
}
