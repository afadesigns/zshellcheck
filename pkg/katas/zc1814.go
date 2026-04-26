// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1814",
		Title:    "Error on `dpkg --force-all` — enables every single `--force-*` option at once",
		Severity: SeverityError,
		Description: "`dpkg --force-all` is shorthand for ~18 distinct `--force-<option>` flags: " +
			"overwrite existing files, install unsigned packages, downgrade, install " +
			"depends-broken, remove essential, and more. The dpkg manual explicitly calls " +
			"this \"almost always a bad idea\". In provisioning scripts it hides the specific " +
			"constraint the author was trying to bypass, and when a later install re-triggers " +
			"the same state the underlying dependency conflict just re-surfaces on the next " +
			"unattended upgrade. Drop `--force-all` and spell out only the `--force-<option>` " +
			"you genuinely need, or fix the upstream conflict.",
		Check: checkZC1814,
	})
}

func checkZC1814(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `dpkg --force-all …` mangles to name=`force-all`.
	if ident.Value == "force-all" {
		return zc1814Hit(cmd)
	}

	if ident.Value != "dpkg" && ident.Value != "apt" && ident.Value != "apt-get" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--force-all" {
			return zc1814Hit(cmd)
		}
		if strings.Contains(v, "Dpkg::Options::=--force-all") {
			return zc1814Hit(cmd)
		}
	}
	return nil
}

func zc1814Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1814",
		Message: "`dpkg --force-all` enables every `--force-*` option at once — " +
			"overwrite, unsigned, downgrade, essential-removal, broken-deps. Drop it " +
			"and spell out only the specific `--force-<option>` you need, or fix the " +
			"upstream conflict.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
