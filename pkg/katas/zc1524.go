package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1524",
		Title:    "Warn on `sysctl -e` / `sysctl -q` — silently skip unknown keys, hide config drift",
		Severity: SeverityWarning,
		Description: "`sysctl -e` and `-q` suppress error output for unknown keys or failed " +
			"writes. That is how a typo in `/etc/sysctl.d/99-hardening.conf` goes unnoticed " +
			"for months — the hardening didn't actually take effect because the key name was " +
			"wrong. Drop `-e`/`-q` in scripts and let errors bubble up; fix the offending " +
			"conffile instead.",
		Check: checkZC1524,
	})
}

func checkZC1524(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sysctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-e" || v == "-q" || v == "-eq" || v == "-qe" {
			return []Violation{{
				KataID: "ZC1524",
				Message: "`sysctl " + v + "` suppresses error output — typos in sysctl.d/ " +
					"conffiles silently skip. Remove and surface the real error.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
