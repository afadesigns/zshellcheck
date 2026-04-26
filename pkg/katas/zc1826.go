// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1826",
		Title:    "Warn on `install -m u+s` / `g+s` — symbolic setuid/setgid bit applied at install time",
		Severity: SeverityWarning,
		Description: "`install -m u+s SRC DEST` (or `g+s` / `ug+s` / `u=rwxs` etc.) applies the " +
			"setuid / setgid bit atomically at copy time — no intermediate `chmod` " +
			"step where a tripwire would fire, no time window where the file exists " +
			"without the special bit. Symbolic forms are easy to miss in review " +
			"because they don't carry the tell-tale leading `4`/`2`/`6` digit that " +
			"numeric-mode detection (see ZC1892) keys off. If DEST is on `$PATH`, " +
			"every local user can invoke the elevated binary. Install setuid / setgid " +
			"binaries only from trusted builds you have reviewed, and prefer narrow " +
			"capabilities (`setcap cap_net_bind_service+ep`) over broad setuid.",
		Check: checkZC1826,
	})
}

func checkZC1826(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "install" {
		return nil
	}
	for i, arg := range cmd.Arguments {
		v := arg.String()
		var mode string
		switch {
		case v == "-m" || v == "--mode":
			if i+1 < len(cmd.Arguments) {
				mode = cmd.Arguments[i+1].String()
			}
		case strings.HasPrefix(v, "-m") && len(v) > 2:
			mode = v[2:]
		case strings.HasPrefix(v, "--mode="):
			mode = strings.TrimPrefix(v, "--mode=")
		}
		if mode == "" {
			continue
		}
		mode = strings.Trim(strings.TrimSpace(mode), "\"'")
		// Numeric setuid / setgid is owned by ZC1892; this kata narrows to
		// symbolic-form setuid/setgid which the numeric scan does not catch.
		if zc1826IsNumericMode(mode) {
			continue
		}
		if zc1826HasSymbolicSetuid(mode) {
			return []Violation{{
				KataID: "ZC1826",
				Message: "`install -m " + mode + "` applies a symbolic setuid/setgid " +
					"bit — easy to miss in review. Use `0755` and grant narrow " +
					"caps with `setcap` instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1826IsNumericMode(mode string) bool {
	if mode == "" {
		return false
	}
	for _, r := range mode {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func zc1826HasSymbolicSetuid(mode string) bool {
	// chmod-style symbolic modes have an `s` or `t` in the perms section.
	// Examples flagged: `u+s`, `g+s`, `ug+s`, `u=rwxs`, `+s`.
	// `s` in the user or group perm slot means setuid / setgid.
	for _, chunk := range strings.Split(mode, ",") {
		if !strings.ContainsAny(chunk, "+=") {
			continue
		}
		if !strings.Contains(chunk, "s") {
			continue
		}
		// Only trip on who-selectors that can carry setuid/setgid (`u`, `g`,
		// `a`, or default/empty `+s` / `=s`) — `o+s` is a no-op.
		idx := strings.IndexAny(chunk, "+=")
		who := chunk[:idx]
		if who == "" || strings.ContainsAny(who, "uga") {
			return true
		}
	}
	return false
}
