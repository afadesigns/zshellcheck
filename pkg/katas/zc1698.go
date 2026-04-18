package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1698",
		Title:    "Warn on `fail2ban-client unban --all` / `stop` — wipes the active brute-force ban list",
		Severity: SeverityWarning,
		Description: "`fail2ban-client unban --all` clears every active ban across every jail; " +
			"`fail2ban-client stop` shuts the service down and flushes its rules. Either " +
			"command restores network access for the exact attacker IPs `fail2ban` has " +
			"already flagged as hostile — usually hundreds of known bots. Target a single " +
			"IP with `fail2ban-client set <jail> unbanip <ip>` or reload a jail with " +
			"`reload <jail>` when you only need to pick up new filter rules.",
		Check: checkZC1698,
	})
}

func checkZC1698(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "fail2ban-client" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	switch cmd.Arguments[0].String() {
	case "stop":
		return zc1698Hit(cmd, "fail2ban-client stop")
	case "unban":
		for _, arg := range cmd.Arguments[1:] {
			if arg.String() == "--all" {
				return zc1698Hit(cmd, "fail2ban-client unban --all")
			}
		}
	}
	return nil
}

func zc1698Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1698",
		Message: "`" + form + "` wipes every active brute-force ban — attacker IPs " +
			"regain access. Target individual IPs with `set <jail> unbanip <ip>` or " +
			"reload a single jail.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
