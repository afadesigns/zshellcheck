package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1436",
		Title:    "`sysctl -w` is ephemeral — persist in `/etc/sysctl.d/*.conf` for surviving reboots",
		Severity: SeverityInfo,
		Description: "`sysctl -w key=value` sets a kernel parameter until the next reboot. For " +
			"configuration that must survive reboots, write a file in `/etc/sysctl.d/` and apply " +
			"with `sysctl --system`. Using only `-w` in provisioning scripts creates silent drift.",
		Check: checkZC1436,
	})
}

func checkZC1436(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sysctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-w" || arg.String() == "--write" {
			return []Violation{{
				KataID: "ZC1436",
				Message: "`sysctl -w` setting is lost on reboot. Persist in `/etc/sysctl.d/*.conf` " +
					"and reload with `sysctl --system`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}
