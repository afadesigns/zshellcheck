// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1702",
		Title:    "Warn on `dpkg-reconfigure` without a noninteractive frontend — hangs in CI",
		Severity: SeverityWarning,
		Description: "`dpkg-reconfigure PACKAGE` opens the package's debconf questions in " +
			"whatever frontend the caller's `DEBIAN_FRONTEND` resolves to — typically a " +
			"terminal dialog that blocks until someone presses a key. Inside a non-" +
			"interactive pipeline (Dockerfile, Ansible task, cloud-init) the call hangs " +
			"until the build times out. Pass `-f noninteractive` (or export " +
			"`DEBIAN_FRONTEND=noninteractive` at the top of the script) and accept the " +
			"debconf defaults; pre-seed any non-default answer with `debconf-set-selections`.",
		Check: checkZC1702,
	})
}

func checkZC1702(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dpkg-reconfigure" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "--frontend=") &&
			strings.Contains(v, "noninteractive") {
			return nil
		}
		if (v == "-f" || v == "--frontend") && i+1 < len(cmd.Arguments) &&
			cmd.Arguments[i+1].String() == "noninteractive" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1702",
		Message: "`dpkg-reconfigure` without `-f noninteractive` opens debconf dialogs — " +
			"non-interactive pipelines hang. Pass `-f noninteractive` or export " +
			"`DEBIAN_FRONTEND=noninteractive`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
