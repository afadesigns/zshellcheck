package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1376",
		Title:    "Avoid `BASH_XTRACEFD` — use Zsh `exec {fd}>file` + `setopt XTRACE`",
		Severity: SeverityWarning,
		Description: "Bash's `BASH_XTRACEFD` redirects `set -x` output to a file descriptor. Zsh " +
			"does not honor this variable; setting it is a silent no-op. To redirect trace output " +
			"in Zsh, open a dedicated fd with `exec {fd}>file` and redirect fd 2 through it: " +
			"`exec 2>&$fd; setopt XTRACE`.",
		Check: checkZC1376,
	})
}

func checkZC1376(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" && ident.Value != "export" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "$BASH_XTRACEFD" || v == "${BASH_XTRACEFD}" ||
			v == "BASH_XTRACEFD" {
			return []Violation{{
				KataID: "ZC1376",
				Message: "`BASH_XTRACEFD` is Bash-only. Zsh ignores it. Redirect trace output " +
					"with `exec {fd}>file; exec 2>&$fd; setopt XTRACE` instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
