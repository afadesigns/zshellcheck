package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1755",
		Title:    "Error on `gcloud sql users {create,set-password} --password PASS` — DB password in argv",
		Severity: SeverityError,
		Description: "`gcloud sql users create USER --instance INST --password PASS` (and the " +
			"`set-password` variant) place the Cloud SQL user password on the command " +
			"line — visible in `ps`, `/proc/<pid>/cmdline`, shell history, and CI logs, " +
			"and stored in Cloud Audit Logs' request payload. Use `--prompt-for-password` " +
			"(interactive) or generate the password server-side in Secret Manager and post " +
			"to the SQL Admin API via `gcloud auth print-access-token` piped to `curl` with " +
			"the body sourced from a file.",
		Check: checkZC1755,
	})
}

func checkZC1755(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gcloud" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "sql" || cmd.Arguments[1].String() != "users" {
		return nil
	}
	sub := cmd.Arguments[2].String()
	if sub != "create" && sub != "set-password" {
		return nil
	}

	prevPwd := false
	for _, arg := range cmd.Arguments[3:] {
		v := arg.String()
		if prevPwd {
			return zc1755Hit(cmd, sub, "--password "+v)
		}
		switch {
		case v == "--password":
			prevPwd = true
		case strings.HasPrefix(v, "--password="):
			return zc1755Hit(cmd, sub, v)
		}
	}
	return nil
}

func zc1755Hit(cmd *ast.SimpleCommand, sub, what string) []Violation {
	return []Violation{{
		KataID: "ZC1755",
		Message: "`gcloud sql users " + sub + " " + what + "` puts the Cloud SQL password " +
			"in argv — visible in `ps`, `/proc`, history, and Cloud Audit Logs. Use " +
			"`--prompt-for-password` or call the SQL Admin API with a body file.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
