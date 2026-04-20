package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1956",
		Title:    "Error on `tailscale up --auth-key=SECRET` — single-use join key visible in argv",
		Severity: SeverityError,
		Description: "`tailscale up --auth-key tskey-auth-…` (and the joined `--auth-key=…` form) " +
			"passes the Tailscale pre-auth key as a command-line argument. Pre-auth keys grant " +
			"full tailnet membership, and short-lived or not, the value ends up in `ps`, " +
			"`/proc/PID/cmdline`, shell history, and any process dump taken before the join " +
			"completes. Read the key from `TS_AUTHKEY` with `tailscale up --authkey-env=TS_AUTHKEY` " +
			"(newer tailscaled), or from a file with `tailscale up --auth-key=file:/etc/ts.key` " +
			"(mode `0400` owned by the provisioning user).",
		Check: checkZC1956,
	})
}

func checkZC1956(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tailscale" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "--auth-key=") || strings.HasPrefix(v, "--authkey=") {
			val := v[strings.IndexByte(v, '=')+1:]
			if zc1956IsLiteralKey(val) {
				return zc1956Hit(cmd, v)
			}
		}
		if (v == "--auth-key" || v == "--authkey") && i+1 < len(cmd.Arguments) {
			val := cmd.Arguments[i+1].String()
			if zc1956IsLiteralKey(val) {
				return zc1956Hit(cmd, v+" "+val)
			}
		}
	}
	return nil
}

func zc1956IsLiteralKey(val string) bool {
	if val == "" {
		return false
	}
	if strings.HasPrefix(val, "file:") {
		return false
	}
	if strings.HasPrefix(val, "$") {
		return false
	}
	return true
}

func zc1956Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1956",
		Message: "`tailscale " + form + "` puts the pre-auth key in argv — visible in " +
			"`ps`/`/proc`/history/crash dumps. Use `--auth-key=file:/etc/ts.key` (mode 0400) " +
			"or `--authkey-env=TS_AUTHKEY`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
