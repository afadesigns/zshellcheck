package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1624",
		Title:    "Error on `az login -p` / `--password` — service-principal secret in process list",
		Severity: SeverityError,
		Description: "`az login -p SECRET` passes the service-principal password as an argv " +
			"element. The expanded value shows up in `ps`, `/proc/<pid>/cmdline`, shell " +
			"history, and audit logs — readable by any local user who can list processes. " +
			"Prefer federated-token OIDC (`--federated-token`), managed identity on the host, " +
			"or interactive device-code flow. If a password is unavoidable, export it as " +
			"`AZURE_PASSWORD` via a protected env var and call plain `az login --service-" +
			"principal` (which reads from env).",
		Check: checkZC1624,
	})
}

func checkZC1624(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "az" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "login" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-p" || v == "--password" {
			return []Violation{{
				KataID: "ZC1624",
				Message: "`az login " + v + "` puts the SP password in argv — visible in " +
					"`ps` / `/proc/<pid>/cmdline`. Use federated-token OIDC, managed " +
					"identity, or `AZURE_PASSWORD` via a protected env var.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
