package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1592",
		Title:    "Warn on `faillock --reset` / `pam_tally2 -r` — clears failed-auth counter",
		Severity: SeverityWarning,
		Description: "Both tools zero the PAM counter that triggers account lockout after too " +
			"many failed logins. A script that resets lockouts — even legitimately, to recover " +
			"locked users — also erases evidence of an ongoing brute-force attempt. Intrusion " +
			"detection relies on those counters for alerting. Do not automate resets; if you " +
			"must, log the prior count and page security on every invocation.",
		Check: checkZC1592,
	})
}

func checkZC1592(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "faillock" && ident.Value != "pam_tally2" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--reset" || v == "-r" {
			return []Violation{{
				KataID: "ZC1592",
				Message: "`" + ident.Value + " " + v + "` clears the PAM failed-auth counter — " +
					"masks ongoing brute force. Log the prior count and alert before resetting.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
