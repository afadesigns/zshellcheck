package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1432",
		Title:    "Warn on `chattr +i` / `+a` — immutable/append-only attributes block cleanup",
		Severity: SeverityWarning,
		Description: "`chattr +i` marks a file immutable (even root cannot delete without " +
			"`chattr -i` first). `+a` makes it append-only. Useful for hardening but often " +
			"surprises later scripts. If you set these in provisioning, document the cleanup " +
			"path; in scripts, verify they are truly needed.",
		Check: checkZC1432,
	})
}

func checkZC1432(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chattr" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		// Detect +i / +a / +u (append, immutable, undelete) flags
		if strings.HasPrefix(v, "+") && (strings.ContainsAny(v, "ia") || v == "+u") {
			return []Violation{{
				KataID: "ZC1432",
				Message: "`chattr +i`/`+a` sets immutable/append-only — blocks later cleanup. " +
					"Document the `-i`/`-a` cleanup path or reconsider if really needed.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
