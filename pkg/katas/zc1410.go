package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1410",
		Title:    "Avoid `bash --rcfile` / `--init-file` — use Zsh `$ZDOTDIR`",
		Severity: SeverityWarning,
		Description: "Bash's `--rcfile FILE` and `--init-file FILE` override the default init " +
			"script. Zsh uses `$ZDOTDIR` to relocate all Zsh rc files (`$ZDOTDIR/.zshrc` etc.). " +
			"Setting ZDOTDIR is cleaner than --rcfile and works for all Zsh init phases.",
		Check: checkZC1410,
	})
}

func checkZC1410(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "bash" && ident.Value != "zsh" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--rcfile" || v == "--init-file" {
			return []Violation{{
				KataID: "ZC1410",
				Message: "`bash --rcfile`/`--init-file` have no Zsh equivalent flag. Use " +
					"`ZDOTDIR=/path zsh` to relocate all Zsh rc files.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
