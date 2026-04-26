// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1864",
		Title:    "Error on `mount -o remount,exec` — re-enables exec on a previously `noexec` mount",
		Severity: SeverityError,
		Description: "Hardened systems mount `/tmp`, `/var/tmp`, `/dev/shm`, and `/home` with " +
			"`noexec` so a dropper cannot chmod and launch a payload out of a world-writable " +
			"directory. `mount -o remount,exec /tmp` (or the narrower `remount,suid`) " +
			"removes that guardrail for the live kernel, and every shell that already had " +
			"`cd /tmp` open picks it up immediately. Most legitimate uses come from install " +
			"scripts that briefly relax `noexec`; those scripts should restore the flag in " +
			"a `trap 'mount -o remount,noexec /tmp' EXIT`. Blanket `remount,exec` without a " +
			"restore path leaves the system in a permanently weakened state until reboot.",
		Check: checkZC1864,
	})
}

func checkZC1864(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mount" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		var opts string
		switch {
		case v == "-o" && i+1 < len(args):
			opts = args[i+1].String()
		case strings.HasPrefix(v, "-o"):
			opts = strings.TrimPrefix(v, "-o")
		default:
			continue
		}
		opts = strings.ToLower(strings.Trim(opts, "\"'"))
		if !strings.Contains(opts, "remount") {
			continue
		}
		if weak := zc1864FirstWeakenedFlag(opts); weak != "" {
			return []Violation{{
				KataID: "ZC1864",
				Message: "`mount -o " + opts + "` re-enables `" + weak + "` on a " +
					"`noexec`/`nosuid`/`nodev`-hardened mount — dropped payloads " +
					"suddenly execute. Pair with a `trap ... EXIT` that restores " +
					"the original flags or skip the remount.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1864FirstWeakenedFlag(opts string) string {
	for _, entry := range strings.Split(opts, ",") {
		entry = strings.TrimSpace(entry)
		switch entry {
		case "exec", "suid", "dev":
			return entry
		}
	}
	return ""
}
