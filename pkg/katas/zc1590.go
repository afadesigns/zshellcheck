package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1590",
		Title:    "Error on `sshpass -p SECRET` — password in process list and history",
		Severity: SeverityError,
		Description: "`sshpass -p SECRET` places the password in argv. It leaks into `ps`, " +
			"`/proc/<pid>/cmdline`, shell history, and audit logs for every process on the box " +
			"that can list processes. The `-f FILE` and `-e` (SSHPASS env) variants keep it off " +
			"argv, but key-based auth is the real fix. Generate an SSH key, authorize it on the " +
			"remote, and drop the password tool entirely.",
		Check: checkZC1590,
	})
}

func checkZC1590(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sshpass" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-p" || (strings.HasPrefix(v, "-p") && len(v) > 2 && v[2] != '=') {
			return []Violation{{
				KataID: "ZC1590",
				Message: "`sshpass -p` places the password in argv — visible in `ps` / " +
					"`/proc/<pid>/cmdline`. Switch to key-based auth, or at least use " +
					"`sshpass -f FILE` / `sshpass -e`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
		if strings.HasPrefix(v, "-p=") {
			return []Violation{{
				KataID: "ZC1590",
				Message: "`sshpass -p` places the password in argv — visible in `ps` / " +
					"`/proc/<pid>/cmdline`. Switch to key-based auth, or at least use " +
					"`sshpass -f FILE` / `sshpass -e`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
