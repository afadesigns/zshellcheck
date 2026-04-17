package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1440",
		Title:    "`usermod -G group user` replaces supplementary groups — use `-aG` to append",
		Severity: SeverityWarning,
		Description: "`usermod -G group user` overwrites the user's supplementary group list — " +
			"any prior group memberships are removed. Users commonly add themselves to `docker` " +
			"or `wheel` via `-G` and inadvertently lose `sudo`/`audio`/other memberships. Always " +
			"pair with `-a` (`-aG`) to append instead of replace.",
		Check: checkZC1440,
	})
}

func checkZC1440(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "usermod" {
		return nil
	}

	hasG := false
	hasA := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-G", "--groups":
			hasG = true
		case "-a", "--append":
			hasA = true
		case "-aG", "-Ga":
			return nil // safe combined flag
		}
	}
	if hasG && !hasA {
		return []Violation{{
			KataID: "ZC1440",
			Message: "`usermod -G` without `-a` overwrites supplementary groups. Use `-aG` to " +
				"append — existing memberships are preserved.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
