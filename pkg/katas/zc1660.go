package katas

import (
	"regexp"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1660ZeroPad = regexp.MustCompile(`%0[1-9][0-9]*d`)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1660",
		Title:    "Style: `printf '%0Nd' $n` — prefer Zsh `${(l:N::0:)n}` left-zero-pad",
		Severity: SeverityStyle,
		Description: "Zero-padding an integer through `printf '%0Nd'` forks a tiny sub-process " +
			"and relies on printf's format-string parser — both things Zsh can avoid. " +
			"`${(l:N::0:)n}` left-pads `$n` with `0` to width N using Zsh parameter " +
			"expansion, no fork, and composes cleanly with other `(q)` / `(L)` / `(U)` " +
			"flags. For right-pad use `${(r:N::0:)n}`; for space padding swap the fill " +
			"character: `${(l:N:)n}` or `${(r:N:)n}`.",
		Check: checkZC1660,
	})
}

func checkZC1660(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "printf" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	fmtArg := cmd.Arguments[0].String()
	if !zc1660ZeroPad.MatchString(fmtArg) {
		return nil
	}

	return []Violation{{
		KataID: "ZC1660",
		Message: "`printf '%0Nd'` forks for zero-padding — prefer Zsh `${(l:N::0:)n}` " +
			"parameter-expansion pad (same for `(r:N::0:)` on the right).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
