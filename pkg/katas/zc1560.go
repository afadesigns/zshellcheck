package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1560",
		Title:    "Error on `pip install --trusted-host` — accepts MITM / plain-HTTP PyPI index",
		Severity: SeverityError,
		Description: "`--trusted-host` tells pip to skip TLS certificate verification for the " +
			"specified host and to allow plain-HTTP URLs from that host. Any MITM on the path " +
			"can substitute packages on install, and a typo in the host name means every " +
			"subsequent `install` from the misspelled host is unauthenticated. Fix the CA " +
			"trust (install the real corporate CA) instead of silencing pip, and keep the " +
			"default `--index-url https://...` over the TLS-verified endpoint.",
		Check: checkZC1560,
	})
}

func checkZC1560(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pip" && ident.Value != "pip3" && ident.Value != "pipx" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--trusted-host" {
			return []Violation{{
				KataID: "ZC1560",
				Message: "`pip --trusted-host` skips TLS verification and allows plain-HTTP " +
					"for that index. Fix the CA trust and keep --index-url on https://.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
