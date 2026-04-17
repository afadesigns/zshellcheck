package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1568",
		Title:    "Error on `useradd -o` / `usermod -o` — allows non-unique UID (alias user)",
		Severity: SeverityError,
		Description: "`-o` (or `--non-unique`) lets `useradd` / `usermod` assign a UID that is " +
			"already in use. The new account has the same kernel identity as the existing one " +
			"but its own login name, password, shell, and home dir. It is indistinguishable in " +
			"`ps` / audit / file ACLs, so a compromise of either account is a compromise of " +
			"both. Pick a fresh UID instead.",
		Check: checkZC1568,
	})
}

func checkZC1568(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "useradd" && ident.Value != "usermod" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-o" || v == "--non-unique" {
			return []Violation{{
				KataID: "ZC1568",
				Message: "`" + ident.Value + " -o` assigns a non-unique UID — the two " +
					"accounts share kernel identity, indistinguishable in audit. Use a " +
					"fresh UID.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
