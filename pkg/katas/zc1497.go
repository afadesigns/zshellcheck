package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1497",
		Title:    "Error on `useradd -u 0` / `usermod -u 0` — creates a second root account",
		Severity: SeverityError,
		Description: "Creating a user with UID 0 makes them a second root — indistinguishable " +
			"from `root` for every access decision, but hiding behind a non-obvious username " +
			"(`backup`, `service`, `svc-updater`). This is a textbook persistence technique. " +
			"If you need privileged but auditable operations, grant sudo rules tied to a " +
			"specific non-0 UID and log via sudo's session plugin.",
		Check: checkZC1497,
	})
}

func checkZC1497(node ast.Node) []Violation {
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

	var prevU bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevU {
			prevU = false
			if v == "0" {
				return zc1497Violation(cmd)
			}
		}
		if v == "-u" || v == "--uid" {
			prevU = true
		}
		if v == "-u0" || v == "--uid=0" {
			return zc1497Violation(cmd)
		}
	}
	return nil
}

func zc1497Violation(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1497",
		Message: "Creating a user with UID 0 produces a second root account — classic " +
			"persistence technique. Use sudo rules tied to a non-0 UID instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
