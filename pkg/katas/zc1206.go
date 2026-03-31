package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1206",
		Title:    "Avoid `crontab -e` in scripts — use `crontab file`",
		Severity: SeverityWarning,
		Description: "`crontab -e` opens an interactive editor which hangs in scripts. " +
			"Use `crontab file` or pipe content with `crontab -` for programmatic cron management.",
		Check: checkZC1206,
	})
}

func checkZC1206(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "crontab" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-e" {
			return []Violation{{
				KataID: "ZC1206",
				Message: "Avoid `crontab -e` in scripts — it opens an interactive editor. " +
					"Use `crontab file` or `echo '...' | crontab -` for programmatic cron management.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
