// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1482",
		Title:    "Error on `docker login -p` / `--password=` — credential in process list",
		Severity: SeverityError,
		Description: "Passing the registry password on the command line puts it in the output of " +
			"`ps`, `/proc/<pid>/cmdline`, and the shell history. On a shared CI runner or a host " +
			"with unprivileged users, that is an immediate leak. Use `--password-stdin` and " +
			"pipe the secret in from `cat /run/secrets/foo` or a credential helper.",
		Check: checkZC1482,
	})
}

func checkZC1482(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" && ident.Value != "buildah" &&
		ident.Value != "skopeo" && ident.Value != "helm" {
		return nil
	}

	// Must see `login` subcommand anywhere.
	var sawLogin bool
	var prevP bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "login" {
			sawLogin = true
			continue
		}
		if !sawLogin {
			continue
		}
		if prevP {
			return zc1482Violation(cmd, "-p "+v)
		}
		if v == "-p" {
			prevP = true
			continue
		}
		if strings.HasPrefix(v, "--password=") {
			return zc1482Violation(cmd, v)
		}
	}
	return nil
}

func zc1482Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1482",
		Message: "`" + what + "` puts the password in ps / /proc / history. Use " +
			"`--password-stdin` piped from a secrets file or credential helper.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
