package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1247",
		Title:    "Avoid `chmod +s` — setuid/setgid bits are security risks",
		Severity: SeverityError,
		Description: "Setting the setuid or setgid bit (`chmod +s` or `chmod u+s`) allows " +
			"files to execute with the owner's privileges, creating privilege escalation risks.",
		Check: checkZC1247,
	})
}

func checkZC1247(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "+s" || val == "u+s" || val == "g+s" || val == "4755" || val == "2755" {
			return []Violation{{
				KataID: "ZC1247",
				Message: "Avoid `chmod +s` — setuid/setgid bits create privilege escalation risks. " +
					"Use `sudo`, capabilities (`setcap`), or dedicated service accounts instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
