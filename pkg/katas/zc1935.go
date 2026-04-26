// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1935",
		Title:    "Warn on `apt autoremove --purge` / `dnf autoremove` — deletes auto-installed deps and their config",
		Severity: SeverityWarning,
		Description: "`apt autoremove --purge` (and `apt-get autoremove --purge`, `dnf autoremove`, " +
			"`zypper rm --clean-deps`) remove every package the resolver thinks is no longer " +
			"required, plus — with `--purge` — their `/etc` config and data dirs. In CI this " +
			"quietly uproots packages someone else installed manually but never `apt-mark " +
			"manual`-ed, and `--purge` makes the removal irreversible. Run a plain `apt " +
			"autoremove --dry-run` in review, mark the keepers with `apt-mark manual`, and " +
			"drop `--purge` from unattended jobs.",
		Check: checkZC1935,
	})
}

func checkZC1935(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var tool string
	switch ident.Value {
	case "apt", "apt-get":
		tool = ident.Value
	case "dnf", "yum":
		tool = ident.Value
	case "zypper":
		tool = "zypper"
	default:
		return nil
	}

	hasAutoremove := false
	hasPurge := false
	hasCleanDeps := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "autoremove":
			hasAutoremove = true
		case "--purge", "--purge-unused":
			hasPurge = true
		case "--clean-deps":
			hasCleanDeps = true
		}
		if v == "rm" && tool == "zypper" {
			hasAutoremove = true
		}
	}
	if !hasAutoremove {
		return nil
	}

	// apt/apt-get autoremove --purge, or dnf/yum autoremove (always purges
	// configs on RPM distros), or zypper rm --clean-deps.
	if (tool == "apt" || tool == "apt-get") && !hasPurge {
		return nil
	}
	if tool == "zypper" && !hasCleanDeps {
		return nil
	}

	return []Violation{{
		KataID: "ZC1935",
		Message: "`" + tool + " autoremove` strips packages the resolver thinks are " +
			"unused plus their configs — uproots packages installed manually but never " +
			"`apt-mark manual`-ed. Dry-run first, mark keepers, drop `--purge` in CI.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
