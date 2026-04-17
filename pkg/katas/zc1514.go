package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1514",
		Title:    "Error on `useradd -p <hash>` / `usermod -p <hash>` — password hash on cmdline",
		Severity: SeverityError,
		Description: "`-p` takes an already-hashed password (crypt(3) format) and writes it " +
			"to `/etc/shadow`. That hash is in `ps`, `/proc/<pid>/cmdline`, and history for as " +
			"long as the process runs — enough time for a co-tenant to grab it and start an " +
			"offline crack. Use `chpasswd` with `--crypt-method=SHA512` reading from stdin, " +
			"or write `/etc/shadow` via a configuration-management tool with proper file " +
			"permissions.",
		Check: checkZC1514,
	})
}

func checkZC1514(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "useradd" && ident.Value != "usermod" && ident.Value != "adduser" {
		return nil
	}

	var prevP bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevP && v != "" && v[0] != '-' {
			return []Violation{{
				KataID: "ZC1514",
				Message: "`" + ident.Value + " -p <hash>` puts the hashed password in ps / " +
					"/proc / history. Use `chpasswd --crypt-method=SHA512` from stdin.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
		prevP = (v == "-p" || v == "--password")
	}
	return nil
}
