package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1844",
		Title:    "Warn on `logger -p local0.info|local7.notice|…` — unreserved facility often uncollected",
		Severity: SeverityWarning,
		Description: "The eight `local0`–`local7` syslog facilities are reserved for site-specific " +
			"use. Most distro `rsyslog` and `systemd-journald` defaults do not route them " +
			"anywhere — they drop on the floor unless someone dropped a matching rule into " +
			"`/etc/rsyslog.d/*.conf`. Scripts that call `logger -p local0.info 'audit: user " +
			"added to wheel'` therefore log to nothing in the audit trail on a stock " +
			"machine. For portable audit-style logging use the POSIX-reserved `auth.notice` " +
			"or `authpriv.info` facility; for application events, pass `-t TAG` and use " +
			"`user.notice` (the default) or a dedicated journald unit.",
		Check: checkZC1844,
	})
}

func checkZC1844(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "logger" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		var facPrio string
		switch {
		case v == "-p" || v == "--priority":
			if i+1 < len(args) {
				facPrio = args[i+1].String()
			}
		case strings.HasPrefix(v, "-p"):
			facPrio = strings.TrimPrefix(v, "-p")
		case strings.HasPrefix(v, "--priority="):
			facPrio = strings.TrimPrefix(v, "--priority=")
		}
		if facPrio == "" {
			continue
		}
		facility := facPrio
		if idx := strings.Index(facPrio, "."); idx >= 0 {
			facility = facPrio[:idx]
		}
		facility = strings.ToLower(strings.Trim(facility, "\"'"))
		if zc1844IsLocalFacility(facility) {
			return []Violation{{
				KataID: "ZC1844",
				Message: "`logger -p " + facPrio + "` writes to a `local*` facility — " +
					"stock `rsyslog`/`journald` rarely collects these. Use " +
					"`auth.notice`/`authpriv.info` for audit events, or " +
					"`user.notice` + `-t TAG` for app logs.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1844IsLocalFacility(f string) bool {
	if len(f) != len("local0") || !strings.HasPrefix(f, "local") {
		return false
	}
	c := f[len("local")]
	return c >= '0' && c <= '7'
}
