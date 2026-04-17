package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1475",
		Title:    "Warn on `setcap` granting dangerous capabilities to a binary (privesc)",
		Severity: SeverityWarning,
		Description: "Adding CAP_SYS_ADMIN, CAP_DAC_OVERRIDE, CAP_DAC_READ_SEARCH, CAP_SYS_PTRACE, " +
			"or CAP_SETUID to a binary lets any user who can execute it perform operations " +
			"roughly equivalent to root — read any file, change any UID, attach ptrace to root " +
			"processes. Scope the capability as narrowly as possible (e.g. CAP_NET_BIND_SERVICE) " +
			"or run the binary under a dedicated service user with a systemd unit.",
		Check: checkZC1475,
	})
}

var setcapDangerous = []string{
	"cap_sys_admin",
	"cap_dac_override",
	"cap_dac_read_search",
	"cap_sys_ptrace",
	"cap_sys_module",
	"cap_setuid",
	"cap_setgid",
	"cap_chown",
}

func checkZC1475(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setcap" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := strings.ToLower(stripOuterQuotes(arg.String()))
		for _, cap := range setcapDangerous {
			if strings.Contains(v, cap) {
				return []Violation{{
					KataID: "ZC1475",
					Message: "`setcap` granting dangerous capability `" + cap + "` makes the " +
						"binary a privesc vector for any executing user. Scope narrower or use a " +
						"dedicated service user.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}
