package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1559",
		Title:    "Warn on `ssh-copy-id -f` / `-o StrictHostKeyChecking=no` — trust-on-first-use key push",
		Severity: SeverityWarning,
		Description: "`ssh-copy-id` opens an SSH connection to deposit the caller's public key. " +
			"With `-f` it overwrites existing `authorized_keys` without prompting; with " +
			"`-o StrictHostKeyChecking=no` it does not verify the host key. Together they " +
			"push a long-term credential at a host the script has never authenticated — a " +
			"network MITM lands a permanent backdoor. Verify the target host's fingerprint " +
			"out of band before pushing keys.",
		Check: checkZC1559,
	})
}

func checkZC1559(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh-copy-id" {
		return nil
	}

	var prevO bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" {
			return zc1559Violation(cmd, "-f")
		}
		if prevO {
			prevO = false
			s := strings.TrimSpace(strings.ToLower(v))
			if s == "stricthostkeychecking=no" || s == "userknownhostsfile=/dev/null" {
				return zc1559Violation(cmd, "-o "+v)
			}
		}
		if v == "-o" {
			prevO = true
		}
	}
	return nil
}

func zc1559Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1559",
		Message: "`ssh-copy-id " + what + "` pushes a long-term credential without host-key " +
			"verification. Verify the fingerprint out of band first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
