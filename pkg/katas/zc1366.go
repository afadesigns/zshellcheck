package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1366",
		Title:    "Use Zsh `limit` instead of POSIX `ulimit` for idiomatic resource queries",
		Severity: SeverityStyle,
		Description: "Zsh provides both `ulimit` (POSIX compatibility) and `limit` (Zsh native). " +
			"`limit` prints human-readable values (`cputime 10 seconds` vs `-t 10`) and accepts " +
			"`unlimited` as a value. Prefer `limit` for Zsh-idiomatic scripts; keep `ulimit` only " +
			"when the script must run under Bash as well.",
		Check: checkZC1366,
	})
}

func checkZC1366(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ulimit" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1366",
		Message: "Use Zsh `limit` (human-readable) or `limit -s` (stdout-only) instead of " +
			"POSIX `ulimit` for Zsh-native resource queries.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
