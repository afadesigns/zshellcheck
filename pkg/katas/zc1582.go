package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1582",
		Title:    "Warn on `bash -x` / `sh -x` / `zsh -x` — traces every command, leaks secrets",
		Severity: SeverityWarning,
		Description: "`-x` turns on xtrace, printing every command (expanded) to stderr before " +
			"it runs. In a CI log that is indexed / shared / archived, any line that touches " +
			"a secret leaks it verbatim — `curl` with a `Bearer` header, `psql` with a " +
			"password, `echo $API_TOKEN > ...`. If you really need tracing, wrap the non-" +
			"secret block with `set -x; ...; set +x` and exclude the secret-handling parts.",
		Check: checkZC1582,
	})
}

func checkZC1582(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "bash" && ident.Value != "sh" && ident.Value != "zsh" &&
		ident.Value != "dash" && ident.Value != "ksh" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-x" || v == "-xv" || v == "-vx" {
			return []Violation{{
				KataID: "ZC1582",
				Message: "`" + ident.Value + " " + v + "` traces every expanded command — CI logs " +
					"leak secrets verbatim. Scope with `set -x; …; set +x`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
