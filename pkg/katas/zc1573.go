package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1573",
		Title:    "Warn on `chattr -i` / `chattr -a` — removes immutable / append-only attribute",
		Severity: SeverityWarning,
		Description: "Removing the immutable (`-i`) or append-only (`-a`) attribute lets the " +
			"file be overwritten or truncated again. When the target is a log file, shadow " +
			"file, or hardened system binary, that flag was explicitly set to make tampering " +
			"noisy. Removing it mid-script is either a one-shot upgrade (follow with the " +
			"`chattr +i` restore) or an anti-forensics step. If it is the former, wrap the " +
			"change in a function and re-set the attribute at the end.",
		Check: checkZC1573,
	})
}

func checkZC1573(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chattr" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-i" || v == "-a" || v == "-ia" || v == "-ai" {
			return []Violation{{
				KataID: "ZC1573",
				Message: "`chattr " + v + "` removes the tamper-evident attribute. If this " +
					"is a one-shot upgrade, re-set the attribute at the end of the block.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
