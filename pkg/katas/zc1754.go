package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1754",
		Title:    "Error on `gh auth status -t` / `--show-token` — prints OAuth token to stdout",
		Severity: SeverityError,
		Description: "`gh auth status -t` (alias `--show-token`) prints the stored GitHub OAuth " +
			"token alongside the status summary. In CI logs, shared terminals, piped to " +
			"`less`/`tee`, or captured via `script`, the token ends up on disk or in " +
			"scrollback where anyone with log access becomes repo-admin. Never combine " +
			"`-t` with `auth status` in automation; if a machine-readable token is needed, " +
			"`gh auth token` prints only the token and makes the secret-handling path " +
			"explicit.",
		Check: checkZC1754,
	})
}

func checkZC1754(node ast.Node) []Violation {
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
	if cmd.Arguments[0].String() != "auth" || cmd.Arguments[1].String() != "status" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if v == "-t" || v == "--show-token" {
			return []Violation{{
				KataID: "ZC1754",
				Message: "`gh auth status " + v + "` prints the OAuth token in the status " +
					"output — CI logs and scrollback become a repo-admin leak. Use " +
					"`gh auth token` in automation so the secret path is explicit.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
