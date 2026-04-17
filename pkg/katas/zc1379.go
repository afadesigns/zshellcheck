package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1379",
		Title:    "Avoid `$PROMPT_COMMAND` — use Zsh `precmd` function",
		Severity: SeverityWarning,
		Description: "Bash runs the command in `$PROMPT_COMMAND` before each prompt. Zsh does not " +
			"honor this variable; the equivalent is a function named `precmd` (or registered via " +
			"`add-zsh-hook precmd name`). Reading `$PROMPT_COMMAND` in Zsh is a no-op.",
		Check: checkZC1379,
	})
}

func checkZC1379(node ast.Node) []Violation {
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
		if strings.Contains(v, "PROMPT_COMMAND") {
			return []Violation{{
				KataID: "ZC1379",
				Message: "`PROMPT_COMMAND` is Bash-only. In Zsh define a `precmd` function or use " +
					"`autoload -Uz add-zsh-hook; add-zsh-hook precmd my_hook`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
