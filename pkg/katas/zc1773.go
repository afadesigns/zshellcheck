package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1773",
		Title:    "Warn on `xargs` without `-r` / `--no-run-if-empty` — runs once on empty input",
		Severity: SeverityWarning,
		Description: "GNU `xargs` (the common default on Linux) invokes the child command once " +
			"with no arguments when its stdin is empty. Paired with a destructive child " +
			"(`xargs rm`, `xargs kill`, `xargs docker stop`) a pipeline that produces zero " +
			"hits silently runs the command with no operand — usually an error at best and a " +
			"footgun at worst. The flag `-r` (GNU) / `--no-run-if-empty` tells xargs to skip " +
			"the call when no items arrive. Add `-r` to every `xargs` pipeline whose producer " +
			"can return no results, or switch to `find ... -exec cmd {} +` which never runs " +
			"the child on empty input. BSD xargs defaults to this behavior, but the portable " +
			"and explicit choice is to pass `-r` and document the intent.",
		Check: checkZC1773,
	})
}

func checkZC1773(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-r" || v == "--no-run-if-empty" {
			return nil
		}
		// Combined short-flag form like `-rt` or `-0r`.
		if len(v) > 1 && v[0] == '-' && v[1] != '-' {
			for _, c := range v[1:] {
				if c == 'r' {
					return nil
				}
			}
		}
	}
	return []Violation{{
		KataID: "ZC1773",
		Message: "`xargs` without `-r` / `--no-run-if-empty` runs the child once with no " +
			"arguments when stdin is empty — a destructive surprise for `xargs rm`, " +
			"`xargs kill`, etc. Add `-r` or switch to `find ... -exec cmd {} +`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
