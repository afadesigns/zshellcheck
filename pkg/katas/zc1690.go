// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"regexp"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var (
	zc1690GitURL = regexp.MustCompile(`^git\+(https?|ssh|file|git)://`)
	zc1690Hash   = regexp.MustCompile(`^[0-9a-f]{7,40}$`)
	zc1690Tag    = regexp.MustCompile(`^v?\d+\.\d+`)
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1690",
		Title:    "Warn on `pip install git+<URL>` without a commit / tag pin",
		Severity: SeverityWarning,
		Description: "`pip install git+https://host/repo[@main]` checks out a moving ref (the " +
			"repository's default branch when no `@` suffix is given, otherwise a branch " +
			"name the attacker can rewrite). Every subsequent install pulls whatever HEAD " +
			"the branch currently points at — no lockfile, no checksum, no reproducibility. " +
			"Pin to a specific commit SHA (`@abc1234…`) or a signed tag (`@v1.2.3`). If a " +
			"proper PyPI release is available, drop the `git+` form entirely and install " +
			"the versioned package.",
		Check: checkZC1690,
	})
}

func checkZC1690(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "pip", "pip3", "pipx", "uv":
	default:
		return nil
	}

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "install" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if !zc1690GitURL.MatchString(v) {
			continue
		}
		at := strings.LastIndex(v, "@")
		// Skip the `@` that's part of `git+ssh://git@host/...` — locate only
		// the refspec `@` that follows the path.
		if at > 0 && at > strings.LastIndex(v, "/") {
			ref := v[at+1:]
			if zc1690Hash.MatchString(ref) || zc1690Tag.MatchString(ref) {
				continue
			}
		}
		return []Violation{{
			KataID: "ZC1690",
			Message: "`" + ident.Value + " install " + v + "` tracks a moving git ref — " +
				"pin to a commit SHA (`@abc1234…`) or signed tag (`@v1.2.3`), or use the " +
				"PyPI release.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
