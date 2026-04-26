// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1817",
		Title:    "Warn on `git push --delete` / `git push -d` / `git push origin :branch` — remote branch removal",
		Severity: SeverityWarning,
		Description: "Deleting a branch on the remote is an irreversible server-side change the " +
			"local reflog cannot rescue. `git push --delete REMOTE BRANCH`, the short `-d`, " +
			"and the legacy `git push REMOTE :BRANCH` colon form all produce the same result: " +
			"the ref vanishes from the server, open pull requests are orphaned, CI runners " +
			"that pinned to the branch lose the target, and recovery needs the last commit " +
			"SHA which may only live in somebody else's local clone. Confirm the remote name, " +
			"check `git branch -r` / `gh pr list --head BRANCH` first, and prefer letting the " +
			"hosting platform delete the branch after a PR merge (with the auto-delete " +
			"setting) rather than scripting the push.",
		Check: checkZC1817,
	})
}

func checkZC1817(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "push" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--delete" || v == "-d" {
			return zc1817Hit(cmd, v)
		}
		// Legacy colon-form: `git push REMOTE :BRANCH` where :BRANCH starts with ":".
		if strings.HasPrefix(v, ":") && len(v) > 1 {
			return zc1817Hit(cmd, "origin "+v)
		}
	}
	return nil
}

func zc1817Hit(cmd *ast.SimpleCommand, flag string) []Violation {
	return []Violation{{
		KataID: "ZC1817",
		Message: "`git push " + flag + "` deletes the remote branch — open PRs are " +
			"orphaned, CI targets disappear, and the last commit SHA can only come " +
			"back from someone else's clone. Let the hosting platform auto-delete " +
			"after merge instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
