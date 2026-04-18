package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1700",
		Title:    "Error on `ldapsearch -w PASSWORD` / `ldapmodify -w PASSWORD` — bind DN password in process list",
		Severity: SeverityError,
		Description: "OpenLDAP client tools (`ldapsearch`, `ldapmodify`, `ldapadd`, `ldapdelete`, " +
			"`ldapmodrdn`, `ldappasswd`, `ldapcompare`) accept the bind password via `-w " +
			"STRING`. Once invoked, the password sits in `/proc/PID/cmdline`, shell " +
			"history, audit records, and any `ps` output — typically granting cn=admin / " +
			"service-account bind over the whole directory. Use `-W` (prompt), `-y " +
			"FILEPATH` (read from a mode-0400 file), or `SASL` auth (`-Y GSSAPI` with " +
			"Kerberos) to keep the secret out of argv.",
		Check: checkZC1700,
	})
}

var zc1700LDAPTools = map[string]struct{}{
	"ldapsearch":  {},
	"ldapmodify":  {},
	"ldapadd":     {},
	"ldapdelete":  {},
	"ldapmodrdn":  {},
	"ldappasswd":  {},
	"ldapcompare": {},
}

func checkZC1700(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if _, ok := zc1700LDAPTools[ident.Value]; !ok {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-w" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		return []Violation{{
			KataID: "ZC1700",
			Message: "`" + ident.Value + " -w PASSWORD` leaks the LDAP bind password into " +
				"`ps` / `/proc/PID/cmdline` — use `-W` to prompt, `-y FILE` for a mode-0400 " +
				"secret file, or SASL (`-Y GSSAPI`).",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}
