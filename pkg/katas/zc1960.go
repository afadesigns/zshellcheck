package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1960",
		Title:    "Warn on `az vm run-command invoke` / `aws ssm send-command` — arbitrary commands on remote VM",
		Severity: SeverityWarning,
		Description: "`az vm run-command invoke --command-id RunShellScript --scripts \"$CMD\"` " +
			"(and the AWS equivalent `aws ssm send-command --document-name AWS-RunShellScript " +
			"--parameters \"commands=['$CMD']\"`) runs arbitrary shell on the target instance " +
			"via the cloud control plane. The identity making the call is whatever role the " +
			"script's credentials carry; if `$CMD` is composed from any operator or attacker " +
			"input, the result is remote code execution through IAM. Gate the call behind " +
			"a shell-escape-safe templater, pin the document version / script to a reviewed " +
			"asset in blob / S3, and require MFA on the invoking role.",
		Check: checkZC1960,
	})
}

func checkZC1960(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "az":
		// Look for `az vm run-command invoke` subcommand sequence.
		if len(cmd.Arguments) >= 3 &&
			cmd.Arguments[0].String() == "vm" &&
			cmd.Arguments[1].String() == "run-command" &&
			(cmd.Arguments[2].String() == "invoke" || cmd.Arguments[2].String() == "create") {
			return zc1960Hit(cmd, "az vm run-command "+cmd.Arguments[2].String())
		}
	case "aws":
		if len(cmd.Arguments) >= 2 &&
			cmd.Arguments[0].String() == "ssm" &&
			cmd.Arguments[1].String() == "send-command" {
			return zc1960Hit(cmd, "aws ssm send-command")
		}
	case "gcloud":
		if len(cmd.Arguments) >= 2 &&
			cmd.Arguments[0].String() == "compute" &&
			cmd.Arguments[1].String() == "ssh" {
			for _, arg := range cmd.Arguments[2:] {
				v := arg.String()
				if v == "--command" || len(v) > 10 && v[:10] == "--command=" {
					return zc1960Hit(cmd, "gcloud compute ssh --command")
				}
			}
		}
	}
	return nil
}

func zc1960Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1960",
		Message: "`" + form + "` runs arbitrary shell on the VM via the cloud control " +
			"plane — operator-composed command strings become IAM-driven RCE. Pin to a " +
			"reviewed asset, template-escape input, require MFA.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
