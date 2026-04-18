package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1696",
		Title:    "Warn on `pnpm install --no-frozen-lockfile` / `yarn install --no-immutable` — CI lockfile drift",
		Severity: SeverityWarning,
		Description: "`pnpm install --no-frozen-lockfile` (pnpm) and `yarn install " +
			"--no-immutable` (yarn 4+) tell the package manager that the lockfile is " +
			"merely a suggestion — any dep resolution change since the lockfile was " +
			"written gets picked up silently. Run that from CI and the artifact no longer " +
			"matches the pinned dependency graph reviewers signed off on. Use `pnpm " +
			"install --frozen-lockfile` (the CI default) or `yarn install --immutable`, " +
			"and let lockfile regen happen only from a dev workstation PR.",
		Check: checkZC1696,
	})
}

func checkZC1696(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pnpm" && ident.Value != "yarn" && ident.Value != "npm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "--no-frozen-lockfile":
			return zc1696Hit(cmd, v)
		case "--no-immutable":
			return zc1696Hit(cmd, v)
		}
	}
	return nil
}

func zc1696Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1696",
		Message: "`" + form + "` allows the lockfile to drift — the CI artifact no " +
			"longer matches the reviewed dependency graph. Use `--frozen-lockfile` / " +
			"`--immutable` in CI.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
