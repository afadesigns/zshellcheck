package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1378",
		Title:    "Avoid uppercase `$DIRSTACK` — Zsh uses lowercase `$dirstack`",
		Severity: SeverityError,
		Description: "Bash's `$DIRSTACK` is the `pushd`/`popd` directory stack. Zsh exposes the " +
			"same stack as lowercase `$dirstack` (per zsh/parameter module). Using uppercase " +
			"`$DIRSTACK` in Zsh accesses an unrelated (and usually empty) variable.",
		Check: checkZC1378,
	})
}

func checkZC1378(node ast.Node) []Violation {
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
		if strings.Contains(v, "DIRSTACK") {
			return []Violation{{
				KataID: "ZC1378",
				Message: "Use lowercase `$dirstack` in Zsh — uppercase `$DIRSTACK` is Bash-only.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
