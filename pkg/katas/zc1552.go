package katas

import (
	"strconv"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1552",
		Title:    "Warn on `openssl dhparam <2048` / `genrsa <2048` — weak key/parameter size",
		Severity: SeverityWarning,
		Description: "Generating DH parameters or RSA keys shorter than 2048 bits is below every " +
			"modern compliance baseline (NIST SP 800-57, BSI TR-02102, Mozilla Server Side TLS). " +
			"A 1024-bit RSA modulus or DH group is within reach of academic precomputation " +
			"(Logjam) and a 512-bit one was broken on commodity hardware in the 1990s. Use " +
			"2048 as a floor and 3072 / 4096 for long-lived keys.",
		Check: checkZC1552,
	})
}

func checkZC1552(node ast.Node) []Violation {
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

	if len(cmd.Arguments) < 2 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "dhparam" && sub != "genrsa" && sub != "gendsa" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n < 2048 {
			return []Violation{{
				KataID: "ZC1552",
				Message: "`openssl " + sub + " " + v + "` uses a weak key/param size — " +
					"modern baselines require 2048+. Use 2048 or 3072/4096 for long-lived keys.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
