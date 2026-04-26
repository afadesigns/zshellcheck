// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1499",
		Title:    "Style: `docker pull <image>` / `:latest` — unpinned image tag",
		Severity: SeverityStyle,
		Description: "Pulling without a tag defaults to `:latest`, which is a moving label. That " +
			"breaks CI reproducibility (yesterday's build passed, today's fails for no reason " +
			"the author changed) and reintroduces supply-chain surface every pull. Pin to a " +
			"specific tag for convenience or to an immutable `@sha256:` digest for production.",
		Check: checkZC1499,
	})
}

func checkZC1499(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" && ident.Value != "nerdctl" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "pull" && sub != "run" {
		return nil
	}

	// Find the image reference — for `pull` it is the first non-flag arg; for
	// `run` it is also typically the first non-flag non-option arg. Bail as
	// soon as we hit an image-looking token.
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if strings.HasPrefix(v, "-") || strings.Contains(v, "=") {
			continue
		}
		// Digest pin — explicit immutable reference.
		if strings.Contains(v, "@sha256:") || strings.Contains(v, "@sha512:") {
			return nil
		}
		// Has a tag — flag only if tag is `latest`.
		if colon := strings.LastIndex(v, ":"); colon != -1 {
			tag := v[colon+1:]
			// Skip port-looking refs like localhost:5000 — not a tag.
			if strings.Contains(tag, "/") {
				// registry-with-port, no tag present → unpinned
				return zc1499Violation(cmd, v)
			}
			if tag == "latest" {
				return zc1499Violation(cmd, v)
			}
			return nil
		}
		// No colon, no @ — bare image means implicit `:latest`.
		return zc1499Violation(cmd, v)
	}
	return nil
}

func zc1499Violation(cmd *ast.SimpleCommand, ref string) []Violation {
	return []Violation{{
		KataID: "ZC1499",
		Message: "`" + ref + "` is unpinned (implicit `:latest`). Pin to a specific tag or " +
			"an immutable `@sha256:` digest for reproducibility.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
