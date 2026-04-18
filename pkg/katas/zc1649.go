package katas

import (
	"strconv"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1649",
		Title:    "Warn on `openssl req -days N` with N > 825 — long-validity certificate",
		Severity: SeverityWarning,
		Description: "CA/Browser Forum capped public TLS cert validity at 825 days in 2018 and " +
			"major browsers tightened it to 398 days in 2020. A cert issued for 3650 days " +
			"(10 years) can not be revoked effectively — once the private key leaks, the " +
			"attacker keeps access until the cert expires naturally. For an internal root CA " +
			"the long validity is defensible; for leaf / server certs keep it under 398 " +
			"days and automate rotation. `-days` over 825 almost always means \"I don't want " +
			"to deal with renewal,\" which is a maintenance smell dressed up as security.",
		Check: checkZC1649,
	})
}

func checkZC1649(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "openssl" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "req" && sub != "x509" && sub != "ca" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-days" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		days, err := strconv.Atoi(cmd.Arguments[i+1].String())
		if err != nil {
			continue
		}
		if days > 825 {
			return []Violation{{
				KataID: "ZC1649",
				Message: "`openssl " + sub + " -days " + cmd.Arguments[i+1].String() +
					"` issues a cert with a long validity. Keep leaf certs under 398 " +
					"days and automate rotation.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
