package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1639AuthHeaders = []string{
	"authorization:",
	"proxy-authorization:",
	"x-api-key:",
	"api-key:",
	"apikey:",
	"x-auth-token:",
	"x-access-token:",
	"cookie:",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1639",
		Title:    "Error on `curl -H 'Authorization: ...'` — credential header in process list",
		Severity: SeverityError,
		Description: "`-H \"Authorization: Bearer $TOKEN\"` (and similar credential-bearing " +
			"headers like `X-Api-Key`, `X-Auth-Token`, `Proxy-Authorization`, `Cookie`) put " +
			"the expanded value in argv. It shows up in `ps`, `/proc/<pid>/cmdline`, shell " +
			"history, and audit logs — every local user who can list processes reads the " +
			"secret. Pass the header via a file with `-H @FILE` or use `--config FILE` so the " +
			"value stays on disk (with 0600 perms), never on the command line.",
		Check: checkZC1639,
	})
}

func checkZC1639(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "curl" && ident.Value != "wget" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v != "-H" && v != "--header" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		header := strings.ToLower(cmd.Arguments[i+1].String())
		for _, h := range zc1639AuthHeaders {
			if strings.Contains(header, h) {
				return []Violation{{
					KataID: "ZC1639",
					Message: "`" + ident.Value + " -H " + cmd.Arguments[i+1].String() +
						"` places the credential in argv — visible via `ps`. Use `-H @FILE`" +
						" or `--config FILE` with 0600 perms.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
	}
	return nil
}
