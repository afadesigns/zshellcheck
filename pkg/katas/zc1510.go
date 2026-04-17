package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1510",
		Title:    "Error on `auditctl -e 0` / `auditctl -D` — disables kernel audit logging",
		Severity: SeverityError,
		Description: "`auditctl -e 0` switches the Linux audit subsystem off, and `auditctl -D` " +
			"deletes every audit rule, including the ones that monitor `/etc/shadow`, `execve`, " +
			"and privilege escalations. Both are textbook anti-forensics steps. If you need to " +
			"temporarily quiet audit for a maintenance window, use `-e 2` (lock enabled + " +
			"immutable) to require a reboot for any further change and document the action.",
		Check: checkZC1510,
	})
}

func checkZC1510(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "auditctl" {
		return nil
	}

	var prevE bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-D" {
			return zc1510Violation(cmd, "-D", "deletes every audit rule")
		}
		if prevE {
			prevE = false
			if v == "0" {
				return zc1510Violation(cmd, "-e 0", "disables audit subsystem")
			}
		}
		if v == "-e" {
			prevE = true
		}
	}
	return nil
}

func zc1510Violation(cmd *ast.SimpleCommand, flag, what string) []Violation {
	return []Violation{{
		KataID: "ZC1510",
		Message: "`auditctl " + flag + "` " + what + " — anti-forensics tactic. Use `-e 2` " +
			"for a reboot-locked maintenance window instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
