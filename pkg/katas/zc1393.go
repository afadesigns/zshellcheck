package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1393",
		Title:    "Avoid `$SRANDOM` — Bash 5.1+ only, read `/dev/urandom` in Zsh",
		Severity: SeverityWarning,
		Description: "Bash 5.1 added `$SRANDOM` as a cryptographically secure 32-bit random value. " +
			"Zsh does not have an equivalent variable. For secure random integers, read bytes " +
			"from `/dev/urandom` (e.g. `(( n = 0x$(od -N4 -An -tx1 /dev/urandom | tr -d ' ') ))`) " +
			"or use an external such as `openssl rand`.",
		Check: checkZC1393,
	})
}

func checkZC1393(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "SRANDOM") {
			return []Violation{{
				KataID: "ZC1393",
				Message: "`$SRANDOM` is Bash 5.1+. In Zsh read `/dev/urandom` directly or use an " +
					"external (`openssl rand`) for secure random integers.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
