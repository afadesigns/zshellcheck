package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1528",
		Title:    "Warn on `chage -M 99999` / `-E -1` — disables password aging / expiry",
		Severity: SeverityWarning,
		Description: "`chage -M 99999` sets the max password age to roughly 273 years (effectively " +
			"never). `chage -E -1` clears the account expiration date. Both silently remove an " +
			"automatic lockout mechanism a compromised credential would otherwise hit. If " +
			"passwords genuinely should not expire (SSO, cert-based auth), encode that in a " +
			"PAM profile rather than per-user `chage`.",
		Check: checkZC1528,
	})
}

func checkZC1528(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chage" {
		return nil
	}

	var prevM, prevE bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevM {
			prevM = false
			if v == "99999" || v == "-1" || v == "0" {
				return zc1528Violation(cmd, "-M "+v)
			}
		}
		if prevE {
			prevE = false
			if v == "-1" {
				return zc1528Violation(cmd, "-E -1")
			}
		}
		switch v {
		case "-M", "--maxdays":
			prevM = true
		case "-E", "--expiredate":
			prevE = true
		}
	}
	return nil
}

func zc1528Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1528",
		Message: "`chage " + what + "` disables password aging — removes automatic lockout. " +
			"Use a PAM profile instead of per-user chage.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
