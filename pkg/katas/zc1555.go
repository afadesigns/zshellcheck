package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1555",
		Title:    "Error on `chmod` / `chown` on `/etc/shadow` or `/etc/sudoers` (managed files)",
		Severity: SeverityError,
		Description: "`/etc/shadow`, `/etc/gshadow`, `/etc/sudoers`, and `/etc/passwd` have " +
			"specific ownership and mode invariants that the distro `passwd`, `chage`, and " +
			"`visudo` tools maintain atomically with file locking. Direct `chmod`/`chown` races " +
			"those tools, can leave the file world-readable mid-modification (leaking the " +
			"shadow file), and will be clobbered on the next `shadow -p` run. Use the proper " +
			"wrapper, or ship a configuration-management drop-in.",
		Check: checkZC1555,
	})
}

func checkZC1555(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chmod" && ident.Value != "chown" && ident.Value != "chgrp" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "/etc/shadow", "/etc/gshadow", "/etc/sudoers", "/etc/passwd",
			"/etc/shadow-", "/etc/gshadow-", "/etc/passwd-", "/etc/sudoers-":
			return []Violation{{
				KataID: "ZC1555",
				Message: "`" + ident.Value + " ... " + v + "` races the distro-managed tool — " +
					"use passwd/chage/visudo or a config-management drop-in.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
