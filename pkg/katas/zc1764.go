package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1764",
		Title:    "Warn on `git commit --no-verify` / `-n` — skips pre-commit and commit-msg hooks",
		Severity: SeverityWarning,
		Description: "`git commit --no-verify` (alias `-n`) bypasses both the pre-commit and " +
			"commit-msg hooks, which are often the last guardrail against leaked secrets, " +
			"formatting drift, or failing tests. The flag is usually a symptom of a hook " +
			"that needs fixing rather than silencing — the exception quickly becomes the " +
			"rule. Fix the blocking hook, carve out a narrow per-file exemption in the " +
			"hook itself, or file a tracked issue, instead of adding `--no-verify` to " +
			"every commit in a script.",
		Check: checkZC1764,
	})
}

func checkZC1764(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "commit" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--no-verify" || v == "-n" {
			return []Violation{{
				KataID: "ZC1764",
				Message: "`git commit " + v + "` skips pre-commit and commit-msg hooks " +
					"— the last guardrail against secret leaks and broken tests. Fix " +
					"the hook or carve a narrow exemption instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
