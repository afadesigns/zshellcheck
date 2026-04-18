package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// mkpasswd flags that take a value as the next argument.
var zc1741ValueFlags = map[string]bool{
	"-m": true, "-S": true, "-R": true, "-P": true,
	"--method": true, "--salt": true, "--rounds": true, "--password-fd": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1741",
		Title:    "Error on `mkpasswd PASSWORD` — clear-text password in process list",
		Severity: SeverityError,
		Description: "`mkpasswd PASSWORD` (whatwg/Debian `whois`-package version) and `mkpasswd " +
			"-m METHOD PASSWORD` hash the password and print the crypt(3) string on stdout. " +
			"Putting PASSWORD on the command line lands it in `ps`, `/proc/<pid>/cmdline`, " +
			"shell history, and the host audit log. Drop the positional password and read " +
			"from stdin (`mkpasswd -s` reads the password from stdin) — pipe the secret " +
			"from a credentials file or vault: `printf %s \"$PASSWORD\" | mkpasswd -s`.",
		Check: checkZC1741,
	})
}

func checkZC1741(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mkpasswd" {
		return nil
	}

	skipNext := false
	hasPositional := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if skipNext {
			skipNext = false
			continue
		}
		if v == "-s" || v == "--stdin" {
			return nil
		}
		if zc1741ValueFlags[v] {
			skipNext = true
			continue
		}
		if v == "" || v[0] == '-' {
			continue
		}
		hasPositional = true
	}

	if !hasPositional {
		return nil
	}

	return []Violation{{
		KataID: "ZC1741",
		Message: "`mkpasswd PASSWORD` puts the cleartext password in argv — visible in " +
			"`ps`, `/proc`, history. Use `mkpasswd -s` and pipe the secret via stdin " +
			"(`printf %s \"$PASSWORD\" | mkpasswd -s`).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
