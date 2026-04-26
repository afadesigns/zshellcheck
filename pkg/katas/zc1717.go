// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1717",
		Title:    "Warn on `docker pull/push --disable-content-trust` — bypasses image signature checks",
		Severity: SeverityWarning,
		Description: "When `DOCKER_CONTENT_TRUST=1` is enforced on a host (or set via `/etc/docker/" +
			"daemon.json`), Docker rejects unsigned image pulls and signs every push. The " +
			"`--disable-content-trust` flag overrides that per command: a `pull` accepts a " +
			"replaced or unsigned image into local storage, a `push` lands an unsigned tag in " +
			"the registry where downstream pulls cannot verify provenance. Drop the flag and " +
			"sign the artifact (`docker trust sign IMAGE:TAG`) instead, or scope the bypass " +
			"with a tight Notary signer policy.",
		Check: checkZC1717,
		Fix:   fixZC1717,
	})
}

// fixZC1717 strips the `--disable-content-trust` flag from a `docker
// {pull,push,build,create,run}` invocation. The argument parses as a
// ConcatenatedExpression whose token literal is just the leading `--`,
// so the whitespace-aware token-strip helper from ZC1238 can't span
// the full literal on its own. Scan the source forward from the
// argument's start offset for the literal flag bytes and delete the
// span (plus the leading whitespace, so the surrounding source stays
// byte-identical).
func fixZC1717(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}
	const flag = "--disable-content-trust"
	sawSub := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if !sawSub {
			switch v {
			case "pull", "push", "build", "create", "run":
				sawSub = true
				continue
			}
		}
		if !sawSub || v != flag {
			continue
		}
		tok := arg.TokenLiteralNode()
		anchor := LineColToByteOffset(source, tok.Line, tok.Column)
		if anchor < 0 {
			return nil
		}
		// The lexer reports `--` as a separate token, then the rest
		// of the flag is concatenated into a ConcatenatedExpression.
		// The token's column lands on the second `-` for the
		// `DASHDASH` token type, so search a 2-byte window around the
		// anchor for the full flag literal in source bytes.
		off := -1
		for _, delta := range []int{-1, 0, 1} {
			cand := anchor + delta
			if cand < 0 || cand+len(flag) > len(source) {
				continue
			}
			if string(source[cand:cand+len(flag)]) == flag {
				off = cand
				break
			}
		}
		if off < 0 {
			return nil
		}
		start := off
		for start > 0 && (source[start-1] == ' ' || source[start-1] == '\t') {
			start--
		}
		end := off + len(flag)
		startLine, startCol := offsetLineColZC1717(source, start)
		if startLine < 0 {
			return nil
		}
		return []FixEdit{{
			Line:    startLine,
			Column:  startCol,
			Length:  end - start,
			Replace: "",
		}}
	}
	return nil
}

func offsetLineColZC1717(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1717(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	var sub string
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if sub == "" {
			switch v {
			case "pull", "push", "build", "create", "run":
				sub = v
				continue
			}
		}
		if sub != "" && v == "--disable-content-trust" {
			return []Violation{{
				KataID: "ZC1717",
				Message: "`docker " + sub + " --disable-content-trust` overrides " +
					"`DOCKER_CONTENT_TRUST=1` — unsigned image moves into the registry " +
					"or local store. Sign the artifact (`docker trust sign`) instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
