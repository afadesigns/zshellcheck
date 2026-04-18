package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1694",
		Title:    "Warn on `ssh -A` / `-o ForwardAgent=yes` — remote host can reuse local keys",
		Severity: SeverityWarning,
		Description: "`ssh -A` (and `-o ForwardAgent=yes`) forwards the caller's `SSH_AUTH_SOCK` " +
			"into the remote session. Anyone with root on the remote (and any process " +
			"that shares its uid) can read the socket and impersonate the caller against " +
			"every host the caller's keys unlock. Prefer `ssh -J JUMP HOST` (ProxyJump) " +
			"for multi-hop access — it keeps the keys on the local side — or configure a " +
			"scoped key for the remote task and copy it in with `ssh-copy-id`. Save key-" +
			"forwarding for interactive use on trusted hosts.",
		Check: checkZC1694,
	})
}

func checkZC1694(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-A" {
			return zc1694Hit(cmd, "-A")
		}
		if v == "-oForwardAgent=yes" {
			return zc1694Hit(cmd, "-o ForwardAgent=yes")
		}
		if v == "-o" && i+1 < len(cmd.Arguments) &&
			cmd.Arguments[i+1].String() == "ForwardAgent=yes" {
			return zc1694Hit(cmd, "-o ForwardAgent=yes")
		}
	}
	return nil
}

func zc1694Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1694",
		Message: "`ssh " + form + "` forwards the caller's `SSH_AUTH_SOCK` into the " +
			"remote — any root on that host can reuse the keys. Use `ssh -J jumphost` " +
			"instead, or a scoped key for the remote task.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
