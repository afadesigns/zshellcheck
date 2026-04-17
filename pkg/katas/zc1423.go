package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1423",
		Title:    "Dangerous: `iptables -F` / `nft flush ruleset` — drops all firewall rules",
		Severity: SeverityWarning,
		Description: "Flushing the firewall ruleset removes every existing rule, typically " +
			"reverting to the default policy. On a remote machine with policy=DROP, this locks " +
			"you out. Save existing rules first (`iptables-save > backup`) and consider " +
			"`iptables-apply` with a rollback timer.",
		Check: checkZC1423,
	})
}

func checkZC1423(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "iptables", "ip6tables":
		for _, arg := range cmd.Arguments {
			if arg.String() == "-F" || arg.String() == "--flush" {
				return []Violation{{
					KataID: "ZC1423",
					Message: "Flushing firewall rules with `-F` removes every rule — risk of " +
						"locking yourself out of remote hosts. Save + use rollback mechanism.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	case "nft":
		for _, arg := range cmd.Arguments {
			if arg.String() == "flush" {
				return []Violation{{
					KataID: "ZC1423",
					Message: "`nft flush ruleset` clears every firewall table — risk of locking " +
						"yourself out of remote hosts. Save + use rollback mechanism.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}

	return nil
}
