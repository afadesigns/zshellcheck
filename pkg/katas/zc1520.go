package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1520",
		Title:    "Warn on `vared <var>` in scripts — reads interactively, hangs non-interactive",
		Severity: SeverityWarning,
		Description: "`vared` is the Zsh interactive line-editor builtin that lets the user edit " +
			"the value of a variable in place. In a non-interactive script (cron job, CI " +
			"runner, ssh-with-command) `vared` has no TTY, so the script either errors out or " +
			"hangs waiting for input that never arrives. For scripted input, read the value " +
			"from stdin (`read varname`), a file, or an environment variable.",
		Check: checkZC1520,
	})
}

func checkZC1520(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "vared" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1520",
		Message: "`vared` requires a TTY — in a non-interactive script it errors or hangs. " +
			"Use `read`, stdin, or environment variables for scripted input.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
