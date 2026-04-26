// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1781GitSubcommands = map[string]bool{
	"clone":     true,
	"fetch":     true,
	"pull":      true,
	"push":      true,
	"ls-remote": true,
	"archive":   true,
}

var zc1781GitUrlSchemes = []string{
	"https://",
	"http://",
	"git+https://",
	"git+http://",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1781",
		Title:    "Error on `git clone https://user:token@host/...` — PAT in argv and git config",
		Severity: SeverityError,
		Description: "A git remote URL in the form `https://user:token@host/path` puts the " +
			"personal access token directly in argv — visible via `ps`, `/proc/PID/cmdline`, " +
			"shell history, and process accounting. `git clone` additionally records the URL " +
			"(including the credentials) in `.git/config` as the `origin` remote, so every " +
			"later `git fetch` / `pull` re-exposes the same token to every user who can read " +
			"that file. Use a credential helper (`git credential-store`, `git credential-" +
			"osxkeychain`), `GIT_ASKPASS` with a secret pulled from an env var, HTTPS + an " +
			"SSH deploy key, or set the token via the `Authorization: Bearer` header with " +
			"`http.extraHeader` from an env-sourced value.",
		Check: checkZC1781,
	})
}

func checkZC1781(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	subIdx := -1
	for i, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			continue
		}
		if zc1781GitSubcommands[v] {
			subIdx = i
		}
		break
	}
	if subIdx == -1 {
		return nil
	}

	for _, arg := range cmd.Arguments[subIdx+1:] {
		v := strings.Trim(arg.String(), "\"'")
		if leak := zc1781HasCredsInURL(v); leak {
			return []Violation{{
				KataID: "ZC1781",
				Message: "`git " + cmd.Arguments[subIdx].String() + " " + v + "` puts " +
					"the token in argv and `.git/config`. Use a credential helper, " +
					"`GIT_ASKPASS`, or `http.extraHeader` with an env-sourced bearer.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1781HasCredsInURL(v string) bool {
	for _, scheme := range zc1781GitUrlSchemes {
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
