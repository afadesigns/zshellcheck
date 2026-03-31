package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1237",
		Title:    "Use `git clean -n` before `git clean -fd`",
		Severity: SeverityWarning,
		Description: "`git clean -fd` permanently deletes untracked files and directories. " +
			"Use `-n` (dry run) first to preview what will be removed.",
		Check: checkZC1237,
	})
}

func checkZC1237(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "clean" {
		return nil
	}

	hasForce := false
	hasDryRun := false

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "-f" || val == "-fd" || val == "-df" || val == "-fdx" {
			hasForce = true
		}
		if val == "-n" {
			hasDryRun = true
		}
	}

	if hasForce && !hasDryRun {
		return []Violation{{
			KataID: "ZC1237",
			Message: "Use `git clean -n` first to preview removals before `git clean -fd`. " +
				"Forced clean permanently deletes untracked files.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
