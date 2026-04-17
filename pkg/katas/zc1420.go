package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1420",
		Title:    "Avoid `chmod +s` / `chmod u+s` — setuid/setgid is a security risk",
		Severity: SeverityWarning,
		Description: "Setuid (mode bit 4000) and setgid (2000) cause the program to run with the " +
			"file-owner's (or group's) privileges, not the caller's. Any bug in such a program " +
			"is a privilege-escalation vector. Reserve setuid for audited, minimal binaries; " +
			"prefer sudo + policy, capabilities, or containers for less-trusted tooling.",
		Check: checkZC1420,
	})
}

func checkZC1420(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := strings.Trim(arg.String(), "'\"")
		// Symbolic setuid/setgid
		if v == "+s" || v == "u+s" || v == "g+s" || v == "+st" || v == "u+st" {
			return []Violation{{
				KataID: "ZC1420",
				Message: "`chmod +s` / `u+s` / `g+s` sets setuid/setgid — privilege-escalation risk. " +
					"Prefer sudo policy, capabilities, or containerization.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
		// Numeric modes starting with 4xxx or 2xxx (setuid/setgid)
		if len(v) == 4 && (v[0] == '4' || v[0] == '2' || v[0] == '6') &&
			v[1] >= '0' && v[1] <= '7' && v[2] >= '0' && v[2] <= '7' && v[3] >= '0' && v[3] <= '7' {
			return []Violation{{
				KataID:  "ZC1420",
				Message: "Numeric mode with leading 4/2/6 sets setuid/setgid — privilege-escalation risk.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	}

	return nil
}
