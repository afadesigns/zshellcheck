package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1656",
		Title:    "Error on `rsync -e 'ssh -o StrictHostKeyChecking=no'` — host-key verify disabled",
		Severity: SeverityError,
		Description: "Disabling host-key verification through rsync's `-e` transport is the " +
			"same attack surface as ZC1479 but easier to miss in review because the ssh flags " +
			"sit inside a quoted string. A MITM on the network path can impersonate the " +
			"remote host and the rsync stream goes straight through. Use `ssh-keyscan` or " +
			"pre-provisioned `~/.ssh/known_hosts` to trust hosts deliberately, and keep " +
			"`StrictHostKeyChecking=yes`.",
		Check: checkZC1656,
	})
}

func checkZC1656(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "rsync" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "StrictHostKeyChecking=no") ||
			strings.Contains(v, "UserKnownHostsFile=/dev/null") {
			return []Violation{{
				KataID: "ZC1656",
				Message: "`rsync -e 'ssh -o StrictHostKeyChecking=no'` disables host-key " +
					"verification — MITM risk. Pre-provision `known_hosts` and keep " +
					"`StrictHostKeyChecking=yes`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
