package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1641",
		Title:    "Error on `kubectl create secret --from-literal=...` / `--docker-password=...`",
		Severity: SeverityError,
		Description: "`kubectl create secret generic --from-literal=KEY=VALUE` and " +
			"`kubectl create secret docker-registry --docker-password=VALUE` put the secret " +
			"content in argv. The expanded value shows up in `ps`, `/proc/<pid>/cmdline`, " +
			"shell history, and audit logs — readable by any local user who can list " +
			"processes. Use `--from-file=KEY=PATH` (reads from a 0600-protected file), " +
			"`--from-env-file=PATH` (reads KEY=VALUE lines), or pipe a manifest into " +
			"`kubectl apply -f -` with base64-encoded `data:` values staged on disk.",
		Check: checkZC1641,
	})
}

func checkZC1641(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "create" || cmd.Arguments[1].String() != "secret" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if strings.HasPrefix(v, "--from-literal=") || strings.HasPrefix(v, "--docker-password=") {
			return []Violation{{
				KataID: "ZC1641",
				Message: "`kubectl create secret " + v + "` puts the secret in argv — " +
					"visible via `ps`. Use `--from-file=KEY=PATH` / `--from-env-file=PATH`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
