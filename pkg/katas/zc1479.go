package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1479",
		Title:    "Error on `ssh/scp -o StrictHostKeyChecking=no` / `UserKnownHostsFile=/dev/null`",
		Severity: SeverityError,
		Description: "Setting `StrictHostKeyChecking=no` or pointing `UserKnownHostsFile` at " +
			"`/dev/null` makes the client accept any server key on the first (and every) " +
			"connection, stripping the protection against MITM that SSH is designed to provide. " +
			"For ephemeral CI targets, pin the host key in `known_hosts` with `ssh-keyscan` and " +
			"verify the fingerprint out of band, or use `StrictHostKeyChecking=accept-new` at " +
			"most.",
		Check: checkZC1479,
	})
}

func checkZC1479(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" && ident.Value != "scp" && ident.Value != "sftp" {
		return nil
	}

	check := func(spec string) []Violation {
		s := strings.TrimSpace(strings.ToLower(spec))
		if s == "stricthostkeychecking=no" {
			return zc1479Violation(cmd, "StrictHostKeyChecking=no")
		}
		if s == "userknownhostsfile=/dev/null" {
			return zc1479Violation(cmd, "UserKnownHostsFile=/dev/null")
		}
		return nil
	}

	var prevO bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevO {
			prevO = false
			if res := check(v); res != nil {
				return res
			}
		}
		if v == "-o" {
			prevO = true
			continue
		}
		if strings.HasPrefix(v, "-o") {
			if res := check(v[2:]); res != nil {
				return res
			}
		}
	}
	return nil
}

func zc1479Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID:  "ZC1479",
		Message: "`" + what + "` disables SSH host-key verification — first MITM owns the connection. Pin the fingerprint in known_hosts instead.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityError,
	}}
}
