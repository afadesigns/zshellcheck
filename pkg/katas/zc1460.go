// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1460",
		Title:    "Warn on `docker run --security-opt seccomp=unconfined` / `apparmor=unconfined`",
		Severity: SeverityWarning,
		Description: "Disabling seccomp or AppArmor removes the syscall / MAC filter that blocks " +
			"most container escape exploits. Only disable these in a known-safe development " +
			"context; production workloads should keep the default profile or ship a stricter " +
			"custom profile.",
		Check: checkZC1460,
	})
}

func checkZC1460(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	var prevOpt bool
	for _, arg := range cmd.Arguments {
		v := arg.String()

		if strings.HasPrefix(v, "--security-opt=") {
			val := strings.TrimPrefix(v, "--security-opt=")
			if isUnconfined(val) {
				return violateZC1460(cmd)
			}
			continue
		}

		if prevOpt {
			prevOpt = false
			if isUnconfined(v) {
				return violateZC1460(cmd)
			}
		}
		if v == "--security-opt" {
			prevOpt = true
		}
	}

	return nil
}

func isUnconfined(v string) bool {
	return v == "seccomp=unconfined" ||
		v == "apparmor=unconfined" ||
		v == "seccomp:unconfined" ||
		v == "apparmor:unconfined"
}

func violateZC1460(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1460",
		Message: "Disabling seccomp or AppArmor removes the main syscall/MAC filter that " +
			"blocks container escapes. Keep the default profile.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
