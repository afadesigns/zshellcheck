package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1903",
		Title:    "Error on `tee /etc/sudoers*` — appends a rule that bypasses `visudo -c` validation",
		Severity: SeverityError,
		Description: "`tee /etc/sudoers` or `tee -a /etc/sudoers.d/<name>` is a common shortcut " +
			"for adding a sudoers rule, but it skips the syntax check that `visudo -c` would " +
			"perform. A malformed line (missing `ALL`, stray colon, unterminated `Cmnd_Alias`) " +
			"makes sudo refuse every invocation — you lock yourself out of root recovery. " +
			"Write the rule to a temporary file, run `visudo -cf /tmp/rule`, and only then " +
			"`install -m 0440 /tmp/rule /etc/sudoers.d/<name>`.",
		Check: checkZC1903,
	})
}

func checkZC1903(node ast.Node) []Violation {
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
		if strings.HasPrefix(v, "-") {
			continue
		}
		if zc1903IsSudoersTarget(v) {
			return []Violation{{
				KataID: "ZC1903",
				Message: "`tee " + v + "` writes a sudoers rule without `visudo -c` " +
					"validation — a syntax error locks every future `sudo` invocation. " +
					"Write to a temp file, run `visudo -cf`, then `install -m 0440` into place.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1903IsSudoersTarget(v string) bool {
	return v == "/etc/sudoers" || strings.HasPrefix(v, "/etc/sudoers.d/")
}
