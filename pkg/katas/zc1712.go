package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1712SecretKeys = []string{
	"password", "passwd", "secret", "token",
	"apikey", "api_key", "api-key",
	"accesskey", "access_key", "access-key",
	"privatekey", "private_key", "private-key",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1712",
		Title:    "Error on `vault kv put PATH password=…` — secret value in process list",
		Severity: SeverityError,
		Description: "`vault kv put PATH key=value` (and the older `vault write PATH key=value`) " +
			"put the value on the command line. When the key name screams secret " +
			"(`password`, `secret`, `token`, `apikey`, `access_key`, `private_key`), the " +
			"cleartext shows up in `ps`, `/proc/<pid>/cmdline`, shell history, and the " +
			"audit log of the calling host — exactly the surface Vault is meant to remove. " +
			"Use `key=@path/to/file` to read from disk, `key=-` to take the value on stdin, " +
			"or `vault kv put -mount=secret PATH @secret.json` for a JSON payload.",
		Check: checkZC1712,
	})
}

func checkZC1712(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "vault" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	var start int
	switch cmd.Arguments[0].String() {
	case "write":
		start = 2
	case "kv":
		if len(cmd.Arguments) < 2 || cmd.Arguments[1].String() != "put" {
			return nil
		}
		start = 3
	default:
		return nil
	}
	if start >= len(cmd.Arguments) {
		return nil
	}

	for _, arg := range cmd.Arguments[start:] {
		v := arg.String()
		eq := strings.IndexByte(v, '=')
		if eq <= 0 {
			continue
		}
		key := strings.ToLower(v[:eq])
		val := v[eq+1:]
		if val == "" || val == "-" || strings.HasPrefix(val, "@") {
			continue
		}
		for _, secret := range zc1712SecretKeys {
			if !strings.Contains(key, secret) {
				continue
			}
			return []Violation{{
				KataID: "ZC1712",
				Message: "`vault " + cmd.Arguments[0].String() + " " + v + "` puts the " +
					"secret value in argv — visible to every local user. Use " +
					"`" + key + "=@FILE` or `" + key + "=-` to read from disk / stdin.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
