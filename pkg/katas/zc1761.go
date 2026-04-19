package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1761",
		Title:    "Warn on `gh gist create --public` — file becomes world-visible and indexed on GitHub",
		Severity: SeverityWarning,
		Description: "`gh gist create --public FILE` (alias `-p`) creates the gist with `public: " +
			"true`. Public gists are listed on `gist.github.com/discover`, crawled by " +
			"search engines, and archived by secondary scrapers — a leaked secret, private " +
			"company snippet, or unreleased note is effectively permanent the moment it " +
			"lands. The default (`public: false`) keeps the gist unlisted and reachable " +
			"only via its URL. Drop `--public` unless public exposure is the explicit " +
			"goal.",
		Check: checkZC1761,
	})
}

func checkZC1761(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "gist" || cmd.Arguments[1].String() != "create" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if v == "--public" || v == "-p" {
			return []Violation{{
				KataID: "ZC1761",
				Message: "`gh gist create " + v + "` publishes the file to the public " +
					"discover feed — search engines crawl it within minutes. Drop the " +
					"flag unless public exposure is the explicit goal.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
