// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1731SecretKeys = []string{
	"password", "passwd", "pwd",
	"secret", "token",
	"apikey", "api_key", "api-key",
	"accesskey", "access_key", "access-key",
	"privatekey", "private_key", "private-key",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1731",
		Title:    "Error on `curl -d 'password=…'` / `wget --post-data='token=…'` — secret in argv",
		Severity: SeverityError,
		Description: "`curl -d` / `--data` / `--data-raw` / `--data-urlencode` and `wget " +
			"--post-data` / `--body-data` put the POST body in argv — visible in `ps`, " +
			"`/proc/<pid>/cmdline`, shell history, and CI logs. When the body contains a " +
			"credential-looking key (`password`, `secret`, `token`, `apikey`, `access_key`, " +
			"`private_key`), the secret leaks the same way an inline `-u user:pass` would. " +
			"Read the value from a file (`curl --data @secret.txt URL`, `--data-binary @-` " +
			"piped from a secrets store) so the secret never reaches the command line.",
		Check: checkZC1731,
	})
}

func checkZC1731(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var dataFlags map[string]bool
	switch ident.Value {
	case "curl":
		dataFlags = map[string]bool{
			"-d":               true,
			"--data":           true,
			"--data-raw":       true,
			"--data-urlencode": true,
			"--data-binary":    true,
		}
	case "wget":
		dataFlags = map[string]bool{
			"--post-data": true,
			"--body-data": true,
		}
	default:
		return nil
	}

	prevFlag := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevFlag != "" {
			if hit := zc1731MatchSecret(v); hit != "" {
				return zc1731Hit(cmd, ident.Value, prevFlag+" "+v, hit)
			}
			prevFlag = ""
			continue
		}
		if dataFlags[v] {
			prevFlag = v
			continue
		}
		// Joined `--data=key=value` form (curl long flags).
		if eq := strings.IndexByte(v, '='); eq > 0 {
			flag := v[:eq]
			if dataFlags[flag] {
				if hit := zc1731MatchSecret(v[eq+1:]); hit != "" {
					return zc1731Hit(cmd, ident.Value, v, hit)
				}
			}
		}
	}
	return nil
}

func zc1731MatchSecret(value string) string {
	body := strings.Trim(value, "'\"")
	if body == "" {
		return ""
	}
	// File reference (curl `@FILE`) or stdin sentinel — safe.
	if body[0] == '@' || body == "-" {
		return ""
	}
	for _, pair := range strings.Split(body, "&") {
		eq := strings.IndexByte(pair, '=')
		if eq <= 0 {
			continue
		}
		key := strings.ToLower(pair[:eq])
		val := pair[eq+1:]
		if val == "" {
			continue
		}
		for _, secret := range zc1731SecretKeys {
			if strings.Contains(key, secret) {
				return key
			}
		}
	}
	return ""
}

func zc1731Hit(cmd *ast.SimpleCommand, tool, flagPart, key string) []Violation {
	return []Violation{{
		KataID: "ZC1731",
		Message: "`" + tool + " " + flagPart + "` puts secret-keyed POST body (`" + key +
			"=…`) in argv — visible in `ps`, `/proc`, history. Read the value from a " +
			"file with `--data @PATH` or `--data-binary @-` piped from a secrets " +
			"store.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
