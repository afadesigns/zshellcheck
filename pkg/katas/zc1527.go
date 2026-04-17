package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1527",
		Title:    "Warn on `crontab -` — replaces cron from stdin, overwrites without diff",
		Severity: SeverityWarning,
		Description: "`crontab -` (or `crontab -u <user> -`) reads a full crontab from stdin and " +
			"replaces the user's existing entries wholesale. Any manual tweak, oncall " +
			"override, or colleague's row is silently deleted. Paired with `curl | crontab -` " +
			"it is a common persistence one-liner. Use `crontab -l > /tmp/old && ... " +
			"crontab -e` with an explicit diff/merge, or ship cron entries via " +
			"`/etc/cron.d/*` managed by config tooling.",
		Check: checkZC1527,
	})
}

func checkZC1527(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "crontab" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-" {
			return []Violation{{
				KataID: "ZC1527",
				Message: "`crontab -` overwrites the user's crontab from stdin — silently " +
					"drops existing rows. Use /etc/cron.d/ files or a diff/merge workflow.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
