// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1251",
		Title:    "Use `mount -o noexec,nosuid` for untrusted media",
		Severity: SeverityWarning,
		Description: "Mounting untrusted filesystems without `noexec,nosuid` allows execution " +
			"of malicious binaries and setuid exploits. Always restrict mount options.",
		Check: checkZC1251,
	})
}

func checkZC1251(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mount" {
		return nil
	}

	hasOptions := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-o" {
			hasOptions = true
		}
	}

	// Only flag mount with device arguments but no -o options
	hasDevice := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 4 && val[:4] == "/dev" {
			hasDevice = true
		}
	}

	if hasDevice && !hasOptions {
		return []Violation{{
			KataID: "ZC1251",
			Message: "Use `mount -o noexec,nosuid,nodev` when mounting external media. " +
				"Without restrictions, mounted filesystems can contain executable exploits.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
