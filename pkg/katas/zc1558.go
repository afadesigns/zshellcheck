// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1558",
		Title:    "Warn on `usermod -aG wheel|sudo|root|adm` — silent privilege group escalation",
		Severity: SeverityWarning,
		Description: "Adding a user to `wheel`, `sudo`, `root`, `adm`, `docker`, or `libvirt` " +
			"from a script grants persistent admin-level access without the review a sudoers " +
			"drop-in or PAM profile would get. `docker` and `libvirt` in particular are " +
			"equivalent to root (spawn privileged containers / raw disk access). Use a " +
			"sudoers.d file scoped to specific commands and audit changes in configuration " +
			"management.",
		Check: checkZC1558,
	})
}

var privGroups = map[string]struct{}{
	"wheel":   {},
	"sudo":    {},
	"root":    {},
	"adm":     {},
	"docker":  {},
	"libvirt": {},
	"lxd":     {},
	"kvm":     {},
	"disk":    {},
}

func checkZC1558(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "usermod" && ident.Value != "gpasswd" && ident.Value != "adduser" {
		return nil
	}

	// Look for the group argument after -aG / -G / -a
	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	for i, a := range args {
		if (a == "-aG" || a == "-Ga" || a == "-G" || a == "--groups" || a == "--append") &&
			i+1 < len(args) {
			groups := strings.Split(args[i+1], ",")
			for _, g := range groups {
				g = strings.TrimSpace(g)
				if _, bad := privGroups[g]; bad {
					return zc1558Violation(cmd, g)
				}
			}
		}
	}
	// gpasswd -a user <group>  => -a then user then group
	if ident.Value == "gpasswd" && len(args) >= 3 && args[0] == "-a" {
		g := args[2]
		if _, bad := privGroups[g]; bad {
			return zc1558Violation(cmd, g)
		}
	}
	return nil
}

func zc1558Violation(cmd *ast.SimpleCommand, group string) []Violation {
	return []Violation{{
		KataID: "ZC1558",
		Message: "Adding user to `" + group + "` grants persistent admin-level access — use a " +
			"scoped sudoers.d drop-in via configuration management.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
