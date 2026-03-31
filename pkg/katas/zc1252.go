package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1252",
		Title:    "Use `getent passwd` instead of `cat /etc/passwd`",
		Severity: SeverityStyle,
		Description: "`cat /etc/passwd` misses users from LDAP, NIS, or SSSD sources. " +
			"`getent passwd` queries NSS and returns all configured user databases.",
		Check: checkZC1252,
	})
}

func checkZC1252(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/etc/passwd" || val == "/etc/group" || val == "/etc/shadow" {
			return []Violation{{
				KataID: "ZC1252",
				Message: "Use `getent` instead of `cat " + val + "`. " +
					"`getent` queries all NSS sources including LDAP and SSSD.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
