package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1566",
		Title:    "Error on `gem install -P NoSecurity|LowSecurity` / `--trust-policy NoSecurity`",
		Severity: SeverityError,
		Description: "RubyGems' trust policy decides what signatures the installer accepts. " +
			"`NoSecurity` skips signature verification entirely; `LowSecurity` warns but still " +
			"installs unsigned gems. On a registry MITM or a hijacked maintainer account those " +
			"policies turn into arbitrary code execution at gem-install time. Use `HighSecurity` " +
			"(reject all but fully-signed) or `MediumSecurity` for hybrid repos.",
		Check: checkZC1566,
	})
}

func checkZC1566(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "gem" {
		return nil
	}

	var prevP bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevP {
			prevP = false
			if v == "NoSecurity" || v == "LowSecurity" {
				return []Violation{{
					KataID: "ZC1566",
					Message: "`gem -P " + v + "` skips signature verification — MITM or " +
						"account compromise becomes RCE at install. Use HighSecurity.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
		if v == "-P" || v == "--trust-policy" {
			prevP = true
		}
	}
	return nil
}
