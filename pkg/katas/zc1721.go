// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1721",
		Title:    "Error on `chmod NNN /dev/<node>` — world-writable device node is local privilege escalation",
		Severity: SeverityError,
		Description: "Granting world-write to a device node hands every local user a primitive: " +
			"`/dev/kvm` becomes a host-root VM-exit gadget, `/dev/uinput` lets any user inject " +
			"keystrokes into the active session, `/dev/loop-control` forges loop devices, " +
			"`/dev/dri/cardN` opens GPU shaders for code-exec, `/dev/mem` / `/dev/kmem` (where " +
			"still permitted) leak kernel state. Keep the kernel-managed default permissions; " +
			"if userspace needs access, add a udev rule that grants it to a specific group, " +
			"never `666` to the world.",
		Check: checkZC1721,
	})
}

func checkZC1721(node ast.Node) []Violation {
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
		if strings.HasPrefix(v, "/dev/") {
			target = v
			continue
		}
		if mode == "" && zc1671WorldWritable(v) {
			mode = v
		}
	}

	if mode == "" || target == "" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1721",
		Message: "`chmod " + mode + " " + target + "` opens a kernel device node to every " +
			"local user — privilege-escalation surface. Use a udev rule that grants the " +
			"specific group access instead of world-write.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
