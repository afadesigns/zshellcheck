// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1543",
		Title:    "Warn on `go install pkg@latest` / `cargo install --git <url>` without rev pin",
		Severity: SeverityWarning,
		Description: "`go install pkg@latest` and `cargo install --git <url>` without `--rev` / " +
			"`--tag` / `--branch` resolve to whatever HEAD is at install time. The next CI " +
			"run can pull a different commit — great for supply-chain attackers to inject " +
			"post-breach, bad for reproducibility. Pin to a specific version tag (`pkg@v1.2.3`) " +
			"or a commit hash (`cargo install --rev abc123 --git ...`).",
		Check: checkZC1543,
	})
}

func checkZC1543(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// go install ...@latest / no @
	if ident.Value == "go" && len(args) >= 2 && args[0] == "install" {
		for _, a := range args[1:] {
			if strings.HasPrefix(a, "-") {
				continue
			}
			if strings.HasSuffix(a, "@latest") || strings.HasSuffix(a, "@master") ||
				strings.HasSuffix(a, "@main") {
				return zc1543Violation(cmd, "go install "+a)
			}
			if !strings.Contains(a, "@") && strings.Contains(a, "/") {
				return zc1543Violation(cmd, "go install "+a+" (no @version)")
			}
		}
	}

	// cargo install --git <url>  with no --rev / --tag / --branch
	if ident.Value == "cargo" && len(args) >= 2 && args[0] == "install" {
		var hasGit, hasPin bool
		for _, a := range args[1:] {
			if a == "--git" {
				hasGit = true
			}
			if a == "--rev" || a == "--tag" || a == "--branch" || a == "--locked" {
				hasPin = true
			}
		}
		if hasGit && !hasPin {
			return zc1543Violation(cmd, "cargo install --git (no --rev/--tag/--branch)")
		}
	}
	return nil
}

func zc1543Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1543",
		Message: "`" + what + "` is unpinned — HEAD-of-default can change between runs. Pin " +
			"to a version tag or commit hash for reproducibility.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
