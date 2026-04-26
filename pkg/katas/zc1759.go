// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1759SecretKeys = []string{
	"password", "passwd",
	"token", "secret",
	"apikey", "api_key", "api-key",
	"accesskey", "access_key", "access-key",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1759",
		Title:    "Error on `vault login TOKEN` / `login -method=… password=…` — credential in process list",
		Severity: SeverityError,
		Description: "Vault accepts credentials on its `login` / `auth` subcommands in two " +
			"argv-leaking shapes: a positional token (`vault login <TOKEN>`) and KEY=VALUE " +
			"pairs for non-token methods (`vault login -method=userpass username=U " +
			"password=P`). Both land the secret in `ps`, `/proc/<pid>/cmdline`, shell " +
			"history, and Vault's audit log request payload. Read the token from stdin " +
			"(`vault login -` with `printf %s \"$TOKEN\" |`) or source `VAULT_TOKEN` from " +
			"a secrets file and run `vault login -method=token`.",
		Check: checkZC1759,
	})
}

func checkZC1759(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "vault" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "login" && sub != "auth" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		// Stdin sentinel, file sentinel, flag-forms are safe.
		if v == "-" || strings.HasPrefix(v, "-") || strings.HasPrefix(v, "@") {
			continue
		}
		// KEY=VALUE pair: flag secret-named keys.
		if eq := strings.IndexByte(v, '='); eq > 0 {
			key := strings.ToLower(v[:eq])
			for _, secret := range zc1759SecretKeys {
				if strings.Contains(key, secret) {
					return zc1759Hit(cmd, sub, v)
				}
			}
			continue
		}
		// Bare positional token.
		return zc1759Hit(cmd, sub, v)
	}
	return nil
}

func zc1759Hit(cmd *ast.SimpleCommand, sub, what string) []Violation {
	return []Violation{{
		KataID: "ZC1759",
		Message: "`vault " + sub + " " + what + "` puts the Vault credential in argv — " +
			"visible in `ps`, `/proc`, history, Vault audit log. Use `vault login -` " +
			"with stdin or source `VAULT_TOKEN` from a secrets file.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
