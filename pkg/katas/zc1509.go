package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1509",
		Title:    "Warn on `trap '' TERM` / `trap - TERM` — ignores/resets fatal signal",
		Severity: SeverityWarning,
		Description: "`trap '' <signal>` makes the signal uninterruptible. `trap - <signal>` " +
			"restores the default disposition, which on `TERM`/`INT`/`HUP` means the script " +
			"exits without running any cleanup handler. Both forms are routinely used to " +
			"harden long-running scripts against accidental `Ctrl-C`, but also to hide from " +
			"`kill` during incident response. Keep the explicit cleanup handler on at least " +
			"`EXIT` so state is always unwound.",
		Check: checkZC1509,
	})
}

func checkZC1509(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "trap" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}
	handler := cmd.Arguments[0].String()
	if handler != "''" && handler != `""` && handler != "-" {
		return nil
	}
	// Signals after the handler.
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		switch v {
		case "TERM", "SIGTERM", "INT", "SIGINT", "HUP", "SIGHUP",
			"QUIT", "SIGQUIT":
			return []Violation{{
				KataID: "ZC1509",
				Message: "`trap " + handler + " " + v + "` silences a fatal signal — cleanup " +
					"handlers never run. Keep at least a cleanup trap on EXIT.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
