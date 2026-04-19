package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// zc1907Weak lists the fs.protected_* + fs.suid_dumpable values that roll back
// kernel safeguards against /tmp-race / hardlink-escalation / FIFO-owner
// symlink-attack patterns.
var zc1907Weak = map[string]string{
	"fs.protected_hardlinks=0": "hardlink following",
	"fs.protected_symlinks=0":  "symlink following in world-writable dirs",
	"fs.protected_fifos=0":     "FIFO open in world-writable dirs",
	"fs.protected_regular=0":   "regular-file open in world-writable dirs",
	"fs.suid_dumpable=1":       "SUID core-dump exposure (1 = group-only)",
	"fs.suid_dumpable=2":       "SUID core-dump exposure (2 = root-readable)",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1907",
		Title:    "Warn on `sysctl -w fs.protected_*=0` / `fs.suid_dumpable=2` — disables /tmp-race safeguards",
		Severity: SeverityWarning,
		Description: "Linux ships `fs.protected_symlinks`, `fs.protected_hardlinks`, " +
			"`fs.protected_fifos`, and `fs.protected_regular` enabled to stop classic " +
			"`/tmp`-race escalation (dangling-symlink, hardlink-pivot, FIFO-open-owner). " +
			"Setting any of them to `0`, or raising `fs.suid_dumpable` above `0`, hands " +
			"unprivileged local users back the primitives. Keep the defaults; if a legacy " +
			"tool genuinely needs them off, scope the change inside a namespace rather than " +
			"flipping the host knob.",
		Check: checkZC1907,
	})
}

func checkZC1907(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sysctl" {
		return nil
	}

	var writing bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-w" || v == "--write" {
			writing = true
			continue
		}
		if strings.HasPrefix(v, "-") {
			continue
		}
		pair := strings.ReplaceAll(v, " ", "")
		if reason, ok := zc1907Weak[pair]; ok && writing {
			return []Violation{{
				KataID: "ZC1907",
				Message: "`sysctl -w " + pair + "` re-enables " + reason + " — classic " +
					"/tmp-race escalation vector. Keep the default; scope any exception in " +
					"a dedicated namespace.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
