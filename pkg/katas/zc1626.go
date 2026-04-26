// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1626SecretKeys = []string{
	"password", "passwd", "secret", "token",
	"apikey", "api_key", "api-key",
	"accesskey", "access_key", "access-key",
	"privatekey", "private_key", "private-key",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1626",
		Title:    "Error on `helm install/upgrade --set KEY=VALUE` with secret-shaped key",
		Severity: SeverityError,
		Description: "`--set` and `--set-string` put the full `KEY=VALUE` pair on the helm " +
			"command line. When the key name looks like a secret (`password`, `secret`, " +
			"`token`, `apikey`, `access_key`, `private_key`), the expanded VALUE appears in " +
			"`ps`, `/proc/<pid>/cmdline`, shell history, and audit logs — readable by any " +
			"local user who can list processes. Put secrets in a protected values file " +
			"(`helm install -f /secure/values.yaml`), or use `--set-file KEY=PATH` so helm " +
			"reads the content from PATH at apply time.",
		Check: checkZC1626,
	})
}

func checkZC1626(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "helm" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "install" && sub != "upgrade" && sub != "template" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v != "--set" && v != "--set-string" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		pair := cmd.Arguments[i+1].String()
		eq := strings.Index(pair, "=")
		if eq < 0 {
			continue
		}
		key := strings.ToLower(pair[:eq])
		for _, s := range zc1626SecretKeys {
			if strings.Contains(key, s) {
				return []Violation{{
					KataID: "ZC1626",
					Message: "`helm " + sub + " " + v + " " + pair + "` places a secret " +
						"value in argv — readable via `ps`. Use `-f values.yaml` or " +
						"`--set-file " + key + "=PATH`.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
	}
	return nil
}
