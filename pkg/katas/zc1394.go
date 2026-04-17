package katas

import (
	"regexp"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// bashVarRE matches `$BASH` used as a standalone variable (not `$BASH_`).
var bashVarRE = regexp.MustCompile(`\$BASH(?:[^_A-Z]|$)`)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1394",
		Title:    "Avoid `$BASH` — Zsh uses `$ZSH_NAME` for the interpreter name",
		Severity: SeverityInfo,
		Description: "Bash's `$BASH` holds the path to the running Bash executable. Zsh's " +
			"equivalent is `$ZSH_NAME` (for the binary name) or `$0` (interactive shell). " +
			"Using `$BASH` in a Zsh script yields empty output.",
		Check: checkZC1394,
	})
}

func checkZC1394(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if bashVarRE.MatchString(v) {
			return []Violation{{
				KataID: "ZC1394",
				Message: "`$BASH` is Bash-only. Zsh exposes the interpreter name via `$ZSH_NAME` " +
					"and the executable path indirectly via `$0`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}
