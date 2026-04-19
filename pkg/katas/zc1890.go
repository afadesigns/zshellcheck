package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1890",
		Title:    "Error on `kadmin -w PASS` / `kinit` with password arg — Kerberos password in argv",
		Severity: SeverityError,
		Description: "`kadmin -w PASS` and `kadmin.local -w PASS` pass the Kerberos admin " +
			"principal's password directly as an argv element. Every `ps`, `/proc/<pid>/" +
			"cmdline`, history file, and CI-pipeline log therefore sees it in plaintext, " +
			"which is catastrophic for an account that can edit the realm's KDC. Use " +
			"`-k -t /etc/krb5.keytab` for non-interactive auth (keytab permissioned to " +
			"root only), or pipe the password through stdin with the `-q` batch form so " +
			"it never rides in argv.",
		Check: checkZC1890,
	})
}

func checkZC1890(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kadmin" && ident.Value != "kadmin.local" && ident.Value != "kpasswd" {
		return nil
	}
	args := cmd.Arguments
	for i := 0; i+1 < len(args); i++ {
		if args[i].String() != "-w" {
			continue
		}
		val := args[i+1].String()
		if val == "" || val[0] == '-' {
			continue
		}
		return []Violation{{
			KataID: "ZC1890",
			Message: "`" + ident.Value + " -w " + val + "` embeds the Kerberos " +
				"admin password in argv — visible to `ps`, `/proc`, shell history. " +
				"Use `-k -t /etc/krb5.keytab` (keytab root-only) or pipe the " +
				"password on stdin.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}
