package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1618",
		Title:    "Warn on `git commit --no-verify` / `git push --no-verify` — bypasses hooks",
		Severity: SeverityWarning,
		Description: "`--no-verify` skips pre-commit, commit-msg, and pre-push hooks. Those " +
			"hooks are where projects run linting, type-checking, unit tests, and secret " +
			"scanning before code lands. A commit or push with `--no-verify` ships code the " +
			"project's own automation would have rejected. Reserve the flag for emergencies " +
			"with a follow-up commit that passes the hooks; scripts should not use it " +
			"routinely.",
		Check: checkZC1618,
	})
}

func checkZC1618(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "commit" && sub != "push" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--no-verify" || (sub == "commit" && v == "-n") {
			return []Violation{{
				KataID: "ZC1618",
				Message: "`git " + sub + " " + v + "` skips pre-" + sub + " / commit-msg " +
					"hooks — lint, test, and secret-scan checks do not run. Reserve for " +
					"emergencies; scripts should pass the hooks.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
