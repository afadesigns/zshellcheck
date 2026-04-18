package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1620",
		Title:    "Error on `tee /etc/sudoers` / `/etc/sudoers.d/*` — writes without `visudo -cf`",
		Severity: SeverityError,
		Description: "`tee` copies stdin to the file with no syntax check. A typo in a sudoers " +
			"rule — a stray comma, a missing `ALL`, an unclosed alias — leaves the file " +
			"unparseable. The next sudo call refuses to load it and on most systems nobody " +
			"can become root until someone boots from rescue media. Pipe the content through " +
			"`visudo -cf /dev/stdin` first, or write to a temp file, validate with " +
			"`visudo -cf`, then atomically `mv` into `/etc/sudoers.d/`.",
		Check: checkZC1620,
	})
}

func checkZC1620(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tee" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "/etc/sudoers" || strings.HasPrefix(v, "/etc/sudoers.d/") {
			return []Violation{{
				KataID: "ZC1620",
				Message: "`tee " + v + "` writes without syntax validation — a typo locks " +
					"everyone out of sudo. Pipe through `visudo -cf /dev/stdin` or stage " +
					"in a temp file and `visudo -cf` before `mv`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
