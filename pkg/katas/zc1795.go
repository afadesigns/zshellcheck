// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1795RemoteActions = map[string]bool{
	"add":     true,
	"set-url": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1795",
		Title:    "Error on `git remote add NAME https://user:token@host/repo` — credentials persisted in `.git/config`",
		Severity: SeverityError,
		Description: "`git remote add NAME URL` and `git remote set-url NAME URL` write the URL " +
			"into `.git/config` verbatim. When the URL embeds a `user:token@host` credential " +
			"segment, every reader of the repo — other local users, a compromised backup, a " +
			"CI cache, or anyone who runs `git config --list` — picks up the secret. It also " +
			"shows up in argv at the moment of creation (visible via `ps` / " +
			"`/proc/PID/cmdline`). Use a credential helper (`git credential-store`, " +
			"`credential-osxkeychain`), `GIT_ASKPASS` sourced from an env var, or HTTPS + a " +
			"deploy SSH key.",
		Check: checkZC1795,
	})
}

func checkZC1795(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "remote" {
		return nil
	}
	if !zc1795RemoteActions[cmd.Arguments[1].String()] {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		v := strings.Trim(arg.String(), "\"'")
		if zc1795UrlHasCreds(v) {
			return []Violation{{
				KataID: "ZC1795",
				Message: "`git remote " + cmd.Arguments[1].String() + " … " + v + "` " +
					"stores the token in `.git/config` and leaks it via argv at " +
					"creation. Use a credential helper, `GIT_ASKPASS`, or an SSH " +
					"deploy key instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1795UrlHasCreds(v string) bool {
	for _, scheme := range []string{"https://", "http://", "git+https://", "git+http://"} {
		if !strings.HasPrefix(v, scheme) {
			continue
		}
		rest := v[len(scheme):]
		at := strings.Index(rest, "@")
		if at <= 0 {
			return false
		}
		userinfo := rest[:at]
		colon := strings.Index(userinfo, ":")
		if colon <= 0 || colon == len(userinfo)-1 {
			return false
		}
		return true
	}
	return false
}
