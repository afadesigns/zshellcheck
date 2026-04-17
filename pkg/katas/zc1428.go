package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1428",
		Title:    "Avoid `curl -u user:pass` — credentials visible in process list",
		Severity: SeverityError,
		Description: "`curl -u user:password` places the credentials in the command line, where " +
			"they show up in `ps`, `/proc/*/cmdline`, shell history, and most audit logs. Use " +
			"`-u user:` with an interactive password prompt, `--netrc`/`--netrc-file` for " +
			"persistent credentials, or a credentials manager.",
		Check: checkZC1428,
	})
}

func checkZC1428(node ast.Node) []Violation {
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

	var sawU bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-u" {
			sawU = true
			continue
		}
		// Next arg after -u containing ':' signals user:pass literal
		if sawU {
			sawU = false
			for i := 0; i < len(v); i++ {
				if v[i] == ':' && i+1 < len(v) && v[i+1] != '\x00' {
					return []Violation{{
						KataID: "ZC1428",
						Message: "`curl -u user:pass` leaks credentials into the process list. " +
							"Use `-u user:` (prompt), `--netrc`, or a credentials manager.",
						Line:   cmd.Token.Line,
						Column: cmd.Token.Column,
						Level:  SeverityError,
					}}
				}
			}
		}
	}

	return nil
}
