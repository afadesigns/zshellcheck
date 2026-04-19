package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1852",
		Title:    "Error on `firewall-cmd --panic-on` — firewalld drops every packet, kills the SSH session",
		Severity: SeverityError,
		Description: "`firewall-cmd --panic-on` puts firewalld into panic mode, which drops every " +
			"inbound and outbound packet regardless of zone or rule. Running this over a " +
			"remote SSH session is the textbook way to lock yourself out: the command " +
			"returns success, the TCP ACK for that reply never arrives, and nobody can " +
			"reach the host until someone visits the console to `--panic-off`. Stage " +
			"panic-mode experiments on a machine you can power-cycle, gate the call behind " +
			"`at now + 5 minutes` with an auto-disable, or use targeted zone rules instead " +
			"of the blanket switch.",
		Check: checkZC1852,
	})
}

func checkZC1852(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `firewall-cmd --panic-on …` mangles the command name to
	// `panic-on` when the flag is the first arg. Cover both shapes.
	if ident.Value == "panic-on" {
		return zc1852Hit(cmd)
	}
	if ident.Value != "firewall-cmd" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "--panic-on" {
			return zc1852Hit(cmd)
		}
	}
	return nil
}

func zc1852Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1852",
		Message: "`firewall-cmd --panic-on` drops every packet regardless of zone — " +
			"an SSH-run call loses the session instantly. Use targeted zone rules; " +
			"if you really need panic mode, gate behind `at now + N minutes … " +
			"firewall-cmd --panic-off`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
