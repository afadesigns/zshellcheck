// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1450",
		Title:    "`pacman -S` / `zypper install` without non-interactive flag hangs in scripts",
		Severity: SeverityWarning,
		Description: "Arch's `pacman -S` waits on confirmation unless `--noconfirm` is passed. " +
			"SUSE's `zypper install` needs `--non-interactive` (or `-n`). Both stall CI pipelines " +
			"and Dockerfiles without these flags.",
		Check: checkZC1450,
	})
}

func checkZC1450(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "pacman":
		hasInstall := false
		hasNoConfirm := false
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if strings.HasPrefix(v, "-S") && !strings.HasPrefix(v, "-Ss") && !strings.HasPrefix(v, "-Si") {
				hasInstall = true
			}
			if v == "--noconfirm" {
				hasNoConfirm = true
			}
		}
		if hasInstall && !hasNoConfirm {
			return []Violation{{
				KataID:  "ZC1450",
				Message: "`pacman -S` without `--noconfirm` hangs in scripts.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	case "zypper":
		hasInstall := false
		hasN := false
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "install" || v == "in" || v == "update" || v == "up" {
				hasInstall = true
			}
			if v == "-n" || v == "--non-interactive" {
				hasN = true
			}
		}
		if hasInstall && !hasN {
			return []Violation{{
				KataID:  "ZC1450",
				Message: "`zypper install` without `--non-interactive` (`-n`) hangs in scripts.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	}

	return nil
}
